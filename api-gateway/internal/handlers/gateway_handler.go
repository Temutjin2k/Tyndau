// internal/handlers/gateway_handler.go
package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"
	"github.com/Temutjin2k/Tyndau/api-gateway/internal/clients"
	"github.com/Temutjin2k/Tyndau/api-gateway/internal/utils"
	pb "github.com/Temutjin2k/TyndauProto/gen/go/user"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type GatewayHandler struct {
	service     *config.Service
	clientPool  *clients.GrpcClientPool
	protoParser *ProtoParser
}

type ProtoParser struct {
	RequestTypeRegistry map[string]proto.Message
}

func NewGatewayHandler(service *config.Service, pool *clients.GrpcClientPool) *GatewayHandler {
	parser := &ProtoParser{
		RequestTypeRegistry: make(map[string]proto.Message),
	}

	// Register inventory service proto types
	parser.RequestTypeRegistry["RegisterRequest"] = &pb.RegisterRequest{}

	return &GatewayHandler{
		service:     service,
		clientPool:  pool,
		protoParser: parser,
	}
}

func (h *GatewayHandler) RegisterRoutes(router *mux.Router) {
	basePath := fmt.Sprintf("/api/%s", h.service.ApiVersion)
	sub := router.PathPrefix(basePath).Subrouter()

	for _, route := range h.service.Routes {
		sub.HandleFunc(route.Path, h.universalHandler(route)).
			Methods(route.Method)
	}
}

func (h *GatewayHandler) universalHandler(route config.RouteConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Get gRPC client connection
		conn, err := h.clientPool.GetConn()
		if err != nil {
			utils.WriteError(w, http.StatusBadGateway, err)
			return
		}

		// Create gRPC request message
		req, err := h.createGRPCRequest(r, route)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		// Invoke gRPC method
		resp, err := h.invokeGRPCMethod(ctx, conn, route, req)
		if err != nil {
			h.handleGRPCError(w, err)
			return
		}

		// Convert response to JSON
		utils.WriteJSON(w, http.StatusOK, resp)
	}
}

func (h *GatewayHandler) createGRPCRequest(r *http.Request, route config.RouteConfig) (proto.Message, error) {
	// Get prototype from registry
	reqProto, ok := h.protoParser.RequestTypeRegistry[route.RequestType]
	if !ok {
		return nil, fmt.Errorf("unknown request type: %s", route.RequestType)
	}

	req := proto.Clone(reqProto)

	// Parse path parameters
	vars := mux.Vars(r)
	for _, param := range route.PathParams {
		if value, exists := vars[strings.ToLower(param)]; exists {
			field := reflect.ValueOf(req).Elem().FieldByName(param)

			if field.IsValid() && field.CanSet() {
				field.SetString(value)
			}
		}
	}

	// Parse query parameters
	query := r.URL.Query()
	for _, param := range route.QueryParams {
		values := query[param]
		if len(values) == 0 {
			continue
		}

		field := reflect.ValueOf(req).Elem().FieldByName(param)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		strValue := values[0]
		err := h.setFieldValue(field, strValue)
		if err != nil {
			return nil, fmt.Errorf("invalid parameter %s: %w", param, err)
		}
	}

	if r.ContentLength > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}

		decoder := protojson.UnmarshalOptions{DiscardUnknown: true}
		if err := decoder.Unmarshal(body, req); err != nil {
			return nil, fmt.Errorf("invalid request body: %w", err)
		}
	}

	return req, nil
}

func (h *GatewayHandler) setFieldValue(field reflect.Value, strValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(strValue)
	case reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

var clientConstructors = map[string]func(*grpc.ClientConn) interface{}{
	"UserService": func(cc *grpc.ClientConn) interface{} {
		return pb.NewAuthClient(cc)
	},
}

func (h *GatewayHandler) invokeGRPCMethod(
	ctx context.Context,
	conn *grpc.ClientConn,
	route config.RouteConfig,
	req proto.Message,
) (proto.Message, error) {
	constructor, ok := clientConstructors[route.GRPCService]
	if !ok {
		return nil, fmt.Errorf("service %s not registered", route.GRPCService)
	}

	client := constructor(conn)

	method := reflect.ValueOf(client).MethodByName(route.GRPCMethod)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found in service %s",
			route.GRPCMethod,
			route.GRPCService)
	}

	params := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
	}

	results := method.Call(params)
	if len(results) != 2 {
		return nil, fmt.Errorf("invalid method signature for %s.%s",
			route.GRPCService,
			route.GRPCMethod)
	}

	// Handle errors
	if err := results[1].Interface(); err != nil {
		return nil, err.(error)
	}

	return results[0].Interface().(proto.Message), nil
}

func (h *GatewayHandler) handleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("%s", st.Message()))
	case codes.InvalidArgument:
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("%s", st.Message()))
	case codes.PermissionDenied:
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("%s", st.Message()))
	case codes.Unauthenticated:
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("%s", st.Message()))
	case codes.DeadlineExceeded:
		utils.WriteError(w, http.StatusGatewayTimeout, fmt.Errorf("request timed out"))
	case codes.ResourceExhausted:
		utils.WriteError(w, http.StatusTooManyRequests, fmt.Errorf("%s", st.Message()))
	default:
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
	}
}

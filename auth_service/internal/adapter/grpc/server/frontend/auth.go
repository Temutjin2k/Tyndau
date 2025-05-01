package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/grpc/server/frontend/dto"
	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	authpb "github.com/Temutjin2k/TyndauProto/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	uc AuthUseCase
}

func NewAuthServer(uc AuthUseCase) *AuthServer {
	return &AuthServer{
		uc: uc,
	}
}

func (h *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := dto.FromRegisterRequest(req)

	newUser, err := h.uc.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		UserId: newUser.ID,
	}, nil
}

func (h *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user := model.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	token, err := h.uc.Login(ctx, user)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Token: token.Token, // or token.Token depending on your `model.Token` struct
	}, nil
}

func (h *AuthServer) IsAdmin(ctx context.Context, req *authpb.IsAdminRequest) (*authpb.IsAdminResponse, error) {
	isAdmin := h.uc.IsAdmin(ctx, req.GetUserId())

	return &authpb.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (h *AuthServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	err := h.uc.Logout(ctx, req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	return &authpb.LogoutResponse{Success: true}, nil
}

package frontend

import (
	"context"
	"fmt"

	"github.com/Temutjin2k/Tyndau/user_service/internal/adapter/grpc/server/frontend/dto"
	userpb "github.com/Temutjin2k/TyndauProto/gen/go/user"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gRPC server for user service
type UserGRPCHandler struct {
	userpb.UnimplementedUserServer
	uc UserUseCase

	log *zerolog.Logger
}

func NewUser(uc UserUseCase, log *zerolog.Logger) *UserGRPCHandler {
	return &UserGRPCHandler{uc: uc, log: log}
}

func (u *UserGRPCHandler) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResonse, error) {
	user := dto.FromCreateRequest(req)

	createdUser, err := u.uc.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.CreateResonse{
		UserId: createdUser.ID,
	}, nil
}

func (u *UserGRPCHandler) ProfileByEmail(ctx context.Context, req *userpb.ProfileByEmailRequest) (*userpb.ProfileResponse, error) {
	// Validate email
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Validate password (if this is a login flow)
	if req.PlainPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, err := u.uc.GetProfileByEmail(ctx, req.Email, req.PlainPassword)
	if err != nil {
		return nil, err
	}

	return dto.ToProfileResponce(user), nil
}

func (u *UserGRPCHandler) Profile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
	if req.UserId < 1 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("wrong ID: %d", req.UserId))
	}

	user, err := u.uc.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return dto.ToProfileResponce(user), nil
}

func (u *UserGRPCHandler) Update(ctx context.Context, req *userpb.UpdateRequest) (*userpb.UpdateResponse, error) {
	user := dto.FromUpdateRequest(req)

	updatedUser, err := u.uc.Update(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return dto.ToUpdateResponce(updatedUser), nil
}

func (u *UserGRPCHandler) Delete(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := u.uc.Delete(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}

package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/grpc/server/frontend/dto"
	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	authpb "github.com/Temutjin2k/TyndauProto/gen/go/auth"
	"github.com/rs/zerolog"
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	uc  AuthUseCase
	log *zerolog.Logger
}

func NewAuthServer(uc AuthUseCase, log *zerolog.Logger) *AuthServer {
	return &AuthServer{
		uc:  uc,
		log: log,
	}
}

func (h *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := dto.FromRegisterRequest(req)

	newUser, err := h.uc.Register(ctx, user)
	if err != nil {
		h.log.Error().Err(err).Msg("Register failed")
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
		h.log.Error().
			Err(err).
			Str("email", user.Email).
			Msg("Login failed")
		return nil, err
	}

	return &authpb.LoginResponse{
		Token: token.Token,
	}, nil
}

func (h *AuthServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponce, error) {
	ok, err := h.uc.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateTokenResponce{
		Succeess: ok,
	}, nil
}

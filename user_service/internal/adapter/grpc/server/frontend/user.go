package frontend

import (
	"context"
	"errors"
	"user_service/internal/adapter/grpc/genproto/userpb"
	"user_service/internal/model"
)

type User struct {
	userpb.UnimplementedUserServiceServer

	uc UserUseCase
}

func NewUser(uc UserUseCase) *User {
	return &User{uc: uc}
}

func (h *User) RegisterUser(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	user := model.User{
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		Password:   req.GetPassword(),
		AvatarLink: req.GetAvatarLink(),
	}

	createdUser, err := h.uc.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		UserId:  createdUser.ID,
		Message: "User registered successfully",
	}, nil
}

func (h *User) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	user := model.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	token, err := h.uc.Authenticate(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userpb.AuthResponse{
		Token: token.Token,
	}, nil
}

func (h *User) GetUserProfile(ctx context.Context, req *userpb.UserID) (*userpb.UserProfile, error) {
	id := req.GetUserId()
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := h.uc.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}

	return &userpb.UserProfile{
		UserId:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		AvatarLink: user.AvatarLink,
		Version:    user.Version,
	}, nil
}

package frontend

import (
	"context"

	userpb "github.com/Temutjin2k/TyndauProto/gen/go/user"
)

type UserGRPCHandler struct {
	userpb.UnimplementedUserServer
	uc UserUseCase
}

func NewUser(uc UserUseCase) *UserGRPCHandler {
	return &UserGRPCHandler{uc: uc}
}

func (u *UserGRPCHandler) Create(context.Context, *userpb.CreateRequest) (*userpb.CreateResonse, error) {
	panic("Implement me")
}

func (u *UserGRPCHandler) Delete(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	panic("implement me")
}

func (u *UserGRPCHandler) Profile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
	panic("implement me")
}

func (u *UserGRPCHandler) Update(context.Context, *userpb.UpdateRequest) (*userpb.UpdateResponse, error) {
	panic("Implement me")
}

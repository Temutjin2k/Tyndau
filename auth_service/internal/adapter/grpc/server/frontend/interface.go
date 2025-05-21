package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
)

type AuthUseCase interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	Login(ctx context.Context, user model.User) (model.Token, error)
	Logout(ctx context.Context, token string) error
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/user_service/internal/model"
)

type UserUseCase interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	GetProfile(ctx context.Context, id int64) (model.User, error)
	GetProfileByEmail(ctx context.Context, email, password string) (model.User, error)
	Delete(ctx context.Context, id int64) error
}

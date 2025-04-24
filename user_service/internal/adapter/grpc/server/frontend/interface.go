package frontend

import (
	"context"
	"user_service/internal/model"
)

type UserUseCase interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	Authenticate(ctx context.Context, user model.User) (model.Token, error)
	GetProfile(ctx context.Context, id int64) (model.User, error)
}

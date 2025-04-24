package usecase

import (
	"context"
	"user_service/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, User model.User) (model.User, error)
	Update(ctx context.Context, update model.User) error
	GetProfile(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id int64) (model.User, error)
}

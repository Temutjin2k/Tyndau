package usecase

import (
	"context"
	"time"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (int64, error)
	User(ctx context.Context, email, password string) (model.User, error)
}

type MailService interface {
	SendWelcome(ctx context.Context, email, name string) error
}

type Producer interface {
	SendEvent(ctx context.Context, subject string, event map[string]any) error
}

type TokenService interface {
	NewToken(user model.User, duration time.Duration) (string, error)
	ValidateToken(token string) (bool, error)
}

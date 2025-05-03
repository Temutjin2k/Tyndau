package usecase

import (
	"context"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (int64, error)
}

type MailService interface {
	SendWelcome(ctx context.Context, email, name string) error
}

type Producer interface {
	SendEvent(ctx context.Context, subject string, event map[string]any) error
}

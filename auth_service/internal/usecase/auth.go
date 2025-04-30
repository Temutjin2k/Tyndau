package usecase

import (
	"context"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	"github.com/rs/zerolog"
)

type AuthUseCase struct {
	userProvider UserService
	mailProvider MailService

	logger *zerolog.Logger
}

func NewAuth(userProvider UserService, mailProvider MailService, logger *zerolog.Logger) *AuthUseCase {
	return &AuthUseCase{
		userProvider: userProvider,
		mailProvider: mailProvider,

		logger: logger,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, user model.User) (model.User, error) {
	userID, err := u.userProvider.Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	user.ID = userID

	err = u.mailProvider.SendWelcome(ctx, user.Email, user.Name)
	if err != nil {
		u.logger.Error().Err(err).Str("email", user.Email).Str("user-name", user.Name).Msg("Failed to send welcome email")
	}

	return user, nil
}

func (u *AuthUseCase) Login(ctx context.Context, user model.User) (model.Token, error) {
	panic("implement me")
}

func (u *AuthUseCase) Logout(ctx context.Context, token string) error {
	panic("implement me")
}

func (u *AuthUseCase) IsAdmin(ctx context.Context, id int64) bool {
	panic("implement me")
}

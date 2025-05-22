package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	"github.com/Temutjin2k/Tyndau/auth_service/pkg/validator"
	"github.com/rs/zerolog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthUseCase struct {
	userProvider  UserService
	mailProvider  MailService
	tokenProvider TokenService

	logger *zerolog.Logger
}

func NewAuth(userProvider UserService, mailProvider MailService, tokenService TokenService, logger *zerolog.Logger) *AuthUseCase {
	return &AuthUseCase{
		userProvider:  userProvider,
		mailProvider:  mailProvider,
		tokenProvider: tokenService,

		logger: logger,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, user model.User) (model.User, error) {
	//Validation
	v := validator.New()

	model.ValidateUser(v, user)
	if ok := v.Valid(); !ok {
		return model.User{}, v
	}

	// Creating new user
	userID, err := u.userProvider.Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	// Assigning id
	user.ID = userID

	// Sends welcome message
	u.logger.Debug().Str("email", user.Email).Str("user-name", user.Name).Msg("Sending email")
	err = u.mailProvider.SendWelcome(ctx, user.Email, user.Name)
	if err != nil {
		u.logger.Error().Err(err).Str("email", user.Email).Str("user-name", user.Name).Msg("Failed to send welcome email")
	}

	return user, nil
}

func (u *AuthUseCase) Login(ctx context.Context, user model.User) (model.Token, error) {
	log := u.logger.With().Str("email", user.Email).Logger()
	log.Info().Msg("attempting to login user")

	//Validation
	v := validator.New()

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(user.Password != "", "password", "must be provided")

	if ok := v.Valid(); !ok {
		log.Error().Err(v).Msg("failed validation")
		return model.Token{}, v
	}

	user, err := u.userProvider.User(ctx, user.Email, user.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return model.Token{}, err
	}

	token, err := u.tokenProvider.NewToken(user, time.Second*5)
	if err != nil {
		return model.Token{}, errors.New("failed to create tokeen")
	}

	return model.Token{
		Token: token,
	}, nil
}

func (u *AuthUseCase) ValidateToken(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, errors.New("empty token provided")
	}

	valid, err := u.tokenProvider.ValidateToken(token)
	if err != nil {
		u.logger.Error().Err(err).Msg("token validation failed")
		return false, ErrInvalidCredentials
	}

	return valid, nil
}

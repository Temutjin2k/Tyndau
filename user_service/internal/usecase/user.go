package usecase

import (
	"context"
	"time"
	"user_service/internal/model"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const jwtSecret = "very-secret"

type UserUseCase struct {
	userRepo UserRepo
}

func NewUser(userRepo UserRepo) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (s *UserUseCase) Register(ctx context.Context, user model.User) (model.User, error) {
	// Check if email is already used
	_, err := s.userRepo.GetProfile(ctx, user.Email)
	if err == nil {
		return model.User{}, model.ErrDuplicateEmail
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user.PasswordHash = string(hashedPassword)

	// Save to database
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

func (s *UserUseCase) Authenticate(ctx context.Context, user model.User) (model.Token, error) {
	// Find user by email
	storedUser, err := s.userRepo.GetProfile(ctx, user.Email)
	if err != nil {
		return model.Token{}, model.ErrAuthenticationFailed
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		return model.Token{}, model.ErrAuthenticationFailed
	}

	// Build JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign it
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return model.Token{}, status.Errorf(codes.Internal, "could not sign token: %v", err)
	}

	return model.Token{Token: signed}, nil
}

func (s *UserUseCase) GetProfile(ctx context.Context, id int64) (model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

package usecase

import (
	"context"

	"github.com/Temutjin2k/Tyndau/user_service/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo UserRepo
}

func NewUser(userRepo UserRepo) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// Create creates new user
func (s *UserUseCase) Create(ctx context.Context, user model.User) (model.User, error) {
	// Check if email is already used
	_, err := s.userRepo.GetProfile(ctx, user.Email)
	if err == nil {
		return model.User{}, model.ErrDuplicateEmail
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
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

// Update updates provided fields
func (s *UserUseCase) Update(ctx context.Context, user model.User) (model.User, error) {
	panic("bro")
}

// GetProfile returns user information
func (s *UserUseCase) GetProfile(ctx context.Context, id int64) (model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// Delete deletes user
func (s *UserUseCase) Delete(ctx context.Context, id int64) {
	panic("bro")
}

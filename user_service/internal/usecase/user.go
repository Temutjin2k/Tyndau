package usecase

import (
	"context"
	"errors"

	"github.com/Temutjin2k/Tyndau/user_service/internal/model"
	"github.com/Temutjin2k/Tyndau/user_service/pkg/def"
	"github.com/Temutjin2k/Tyndau/user_service/pkg/validator"
	"github.com/rs/zerolog"
)

type UserUseCase struct {
	userRepo UserRepo

	log *zerolog.Logger
}

func NewUser(userRepo UserRepo, log *zerolog.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		log:      log,
	}
}

// Create creates new user
func (s *UserUseCase) Create(ctx context.Context, user model.User) (model.User, error) {

	v := validator.New()
	if model.ValidateUser(v, user); !v.Valid() {
		return model.User{}, v
	}

	// Check if email is already used
	_, err := s.userRepo.GetProfile(ctx, user.Email)
	if err == nil {
		return model.User{}, model.ErrDuplicateEmail
	}

	// Hash the password
	err = user.HashPassword()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to hash password")
		return model.User{}, errors.New("failed to hash password")
	}
	// Save to database
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return model.User{}, model.ErrCreateUser
	}

	return createdUser, nil
}

// Update updates provided fields
func (s *UserUseCase) Update(ctx context.Context, user model.User) (model.User, error) {
	dbUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return model.User{}, model.ErrNotFound
	}

	if user.Version != dbUser.Version {
		return model.User{}, model.ErrEditConflict
	}

	// TODO update password
	updatedUser := model.UserUpdate{
		Name:       def.Pointer(user.Name),
		AvatarLink: def.Pointer(user.AvatarLink),
	}

	// If the updatedUser.Name value is nil then we know that no corresponding "name" key/
	// value pair was provided in the request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new name
	// value. Importantly, because updatedUser.Name is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.
	if updatedUser.Name != nil {
		dbUser.Name = *updatedUser.Name
	}

	if updatedUser.AvatarLink != nil {
		dbUser.AvatarLink = *updatedUser.AvatarLink
	}

	// Validating updated user
	v := validator.New()
	if model.ValidateUpdatedUser(v, dbUser); !v.Valid() {
		return model.User{}, v
	}

	// Updating user
	err = s.userRepo.Update(ctx, &dbUser)
	if err != nil {
		return model.User{}, err
	}

	return dbUser, nil
}

// GetProfile returns user information
func (s *UserUseCase) GetProfile(ctx context.Context, id int64) (model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, model.ErrNotFound
	}

	return user, nil
}

// Delete deletes user
func (s *UserUseCase) Delete(ctx context.Context, id int64) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return model.ErrDeleteUser
	}

	return nil
}

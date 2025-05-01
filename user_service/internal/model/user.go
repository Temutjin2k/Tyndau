package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID           int64
		CreatedAt    time.Time
		Name         string
		Email        string
		AvatarLink   string
		Password     string
		PasswordHash string
		Version      int32

		IsDeleted bool
	}

	UserUpdate struct {
		ID           *int64
		CreatedAt    *time.Time
		Name         *string
		Email        *string
		AvatarLink   *string
		Password     *string
		PasswordHash *string
		Version      *int32

		IsDeleted *bool
	}

	Token struct {
		Token string
	}
)

// HashPassword hashes the Password field using bcrypt.GenerateFromPassword and sets it to PasswordHash field
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)

	return nil
}

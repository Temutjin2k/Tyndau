package model

import "time"

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
		IsDeleted    bool
	}

	Token struct {
		Token string
	}
)

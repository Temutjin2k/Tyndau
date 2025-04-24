package dao

import (
	"time"
	"user_service/internal/model"
)

type User struct {
	ID           int64     `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	AvatarLink   string    `db:"avatar_link"`
	PasswordHash []byte    `db:"password_hash"`
	Version      int32     `db:"version"`
	IsDeleted    bool      `db:"is_deleted"`
}

func FromUser(m model.User) User {
	return User{
		ID:           m.ID,
		CreatedAt:    m.CreatedAt,
		Name:         m.Name,
		Email:        m.Email,
		AvatarLink:   m.AvatarLink,
		PasswordHash: []byte(m.PasswordHash),
		Version:      m.Version,
		IsDeleted:    m.IsDeleted,
	}
}

func ToUser(u User) model.User {
	return model.User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		Name:         u.Name,
		Email:        u.Email,
		AvatarLink:   u.AvatarLink,
		PasswordHash: string(u.PasswordHash),
		Version:      u.Version,
		IsDeleted:    u.IsDeleted,
	}
}

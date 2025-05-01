package model

import "errors"

var (
	ErrNotFound             = errors.New("requested resource not found")
	ErrEditConflict         = errors.New("unable to update the record due to an edit conflict")
	ErrDuplicateEmail       = errors.New("user already exists")
	ErrAuthenticationFailed = errors.New("invalid credentials")
	ErrDeleteUser           = errors.New("failed to delete user")
	ErrCreateUser           = errors.New("failed to create user")
)

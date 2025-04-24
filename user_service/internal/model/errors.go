package model

import "errors"

var (
	ErrDuplicateEmail       = errors.New("user already exists")
	ErrAuthenticationFailed = errors.New("invalid credentials")
)

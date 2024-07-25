package storage

import "errors"

var (
	ErrUserExists      = errors.New("user exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrMessageNotFound = errors.New("message not found")
)

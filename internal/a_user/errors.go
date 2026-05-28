package a_user

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken = errors.New("invalid token")
	ErrLoginOrPasswordIncorrect = errors.New("login or password is incorrect")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrPasswordTooLong = errors.New("password is too long")
)
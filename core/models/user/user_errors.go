package user

import "errors"

var (
	ErrDublicateUser = errors.New("User does already exists")
)

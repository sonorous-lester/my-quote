package exceptions

import "errors"

var (
	ValidInput            = errors.New("invalid input")
	ValidEmailAddr        = errors.New("invalid email address")
	InvalidPasswordLength = errors.New("password length should be 6-15 characters")
)

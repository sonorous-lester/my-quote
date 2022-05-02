package exceptions

import "errors"

var (
	InvalidInput          = errors.New("invalid input")
	InvalidEmailAddr      = errors.New("invalid email address")
	InvalidPasswordLength = errors.New("password length should be 6-15 characters")
)

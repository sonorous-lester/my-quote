package exceptions

import "errors"

var (
	InvalidInput          = errors.New("invalid input")
	InvalidEmailAddr      = errors.New("invalid email address")
	InvalidPasswordLength = errors.New("password length should be 6-15 characters")
	AuthError             = errors.New("email or password incorrect")
	UserNotExists         = errors.New("user not exists")
	UserExists            = errors.New("user exists")
	ServerError           = errors.New("server error")
)

package auth

import "myquote/domain/models"

type Usecase interface {
	Register(user NewUser) error
	Login(a Anonymous) (models.User, error)
}

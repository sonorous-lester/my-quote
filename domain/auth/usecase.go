package auth

import "myquote/domain/models"

type Usecase interface {
	Register(user NewUser) error
	Login(i LoginInfo) (models.User, error)
}

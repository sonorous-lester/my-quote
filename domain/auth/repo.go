package auth

import "myquote/domain/models"

type Repository interface {
	FindUser(email string) (bool, models.UserModel, error)
	Register(name string, email string, password string) error
	UpdateToken(user models.UserModel, token string) error
}

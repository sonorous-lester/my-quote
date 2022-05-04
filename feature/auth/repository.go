package auth

import (
	"errors"
	"gorm.io/gorm"
	"myquote/domain"
	"myquote/domain/models"
)

type Repository struct {
	l  domain.Logger
	db *gorm.DB
}

func (r *Repository) FindUser(email string) (bool, models.UserModel, error) {
	var user models.UserModel
	result := r.db.Table("users").First(&user, "email = ?", email)

	userNotFound := result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound)
	if userNotFound {
		return false, models.UserModel{}, nil
	}
	sqlServerError := result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound)
	if sqlServerError {
		r.l.Debugf("find user error, user eamil is :%s.\n The error message is: %s", email, result.Error.Error())
		return false, models.UserModel{}, result.Error
	}

	return true, user, nil
}

func (r *Repository) Register(name string, email string, password string) error {
	user := models.UserModel{Name: name, Email: email, Hashed: password}
	result := r.db.Create(&user)
	if result.Error != nil {
		r.l.Debugf("create user error; username: %s, email: %s\n The error message: %s", name, email, result.Error.Error())
		return result.Error
	}
	return nil
}

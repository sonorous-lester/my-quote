package auth

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"myquote/domain"
	"myquote/domain/models"
)

type Repository struct {
	l  domain.Logger
	db *gorm.DB
}

func (r *Repository) FindUser(email string) (bool, error) {
	var user models.UserModel
	result := r.db.First(&user, "email = ?", email)

	userNotFound := result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound)
	if userNotFound {
		return false, nil
	}
	sqlServerError := result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound)
	if sqlServerError {
		logrus.Debugf("find user error, user eamil is :%s.\n The error message is: %s", email, result.Error.Error())
		return false, result.Error
	}

	return true, nil
}

func (r *Repository) Register(name string, email string, password string) error {
	user := models.UserModel{Name: name, Email: email, Password: password}
	result := r.db.Create(&user)
	if result.Error != nil {
		r.l.Debugf("create user error; username: %s, email: %s\n The error message: %s", name, email, result.Error.Error())
		return result.Error
	}
	return nil
}

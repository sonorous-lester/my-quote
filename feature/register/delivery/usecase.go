package delivery

import (
	"myquote/domain"
	"myquote/domain/common"
	"myquote/domain/exceptions"
	"myquote/domain/register"
)

type RegisterUsecase struct {
	l  domain.Logger
	r  register.Repository
	pv common.Validator
	ev common.Validator
}

func NewRegisterUsecase(logger domain.Logger, repository register.Repository, passwordValidator common.Validator, emailValidator common.Validator) *RegisterUsecase {
	return &RegisterUsecase{
		l:  logger,
		r:  repository,
		pv: passwordValidator,
		ev: emailValidator,
	}
}

func (uc *RegisterUsecase) Register(user register.NewUser) error {
	if !uc.ev.Validate(user.Email) {
		uc.l.Debugf("invalid email addr: %s", user.Email)
		return exceptions.InvalidEmailAddr
	}
	if !uc.pv.Validate(user.Password) {
		uc.l.Debugf("invalid password length: %d", len(user.Password))
		return exceptions.InvalidPasswordLength
	}

	find, _ := uc.r.FindUser(user.Email)
	if find {
		return exceptions.UserExists
	}
	return nil
}

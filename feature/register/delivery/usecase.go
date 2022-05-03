package delivery

import (
	"myquote/domain"
	"myquote/domain/common"
	"myquote/domain/exceptions"
	"myquote/domain/register"
)

type RegisterUsecase struct {
	l     domain.Logger
	r     register.Repository
	pv    common.Validator
	ev    common.Validator
	hashv common.HashValidator
}

func NewRegisterUsecase(logger domain.Logger, repository register.Repository, passwordValidator common.Validator, emailValidator common.Validator, hashValidator common.HashValidator) *RegisterUsecase {
	return &RegisterUsecase{
		l:     logger,
		r:     repository,
		pv:    passwordValidator,
		ev:    emailValidator,
		hashv: hashValidator,
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

	find, err := uc.r.FindUser(user.Email)
	if find {
		return exceptions.UserExists
	}
	if err != nil {
		return exceptions.ServerError
	}

	hash, err := uc.hashv.Hash(user.Password)
	if err != nil {
		uc.l.Debugf("hash password error: %s", err.Error())
		return exceptions.ServerError
	}

	err = uc.r.Register(user.Email, hash)
	if err != nil {
		return exceptions.ServerError
	}

	return nil
}

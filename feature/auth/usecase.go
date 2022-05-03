package auth

import (
	"myquote/domain"
	"myquote/domain/auth"
	"myquote/domain/common"
	"myquote/domain/exceptions"
)

type Usecase struct {
	l     domain.Logger
	r     auth.Repository
	pv    common.Validator
	ev    common.Validator
	hashv common.HashValidator
}

func NewUsecase(logger domain.Logger, repository auth.Repository, passwordValidator common.Validator, emailValidator common.Validator, hashValidator common.HashValidator) *Usecase {
	return &Usecase{
		l:     logger,
		r:     repository,
		pv:    passwordValidator,
		ev:    emailValidator,
		hashv: hashValidator,
	}
}

func (uc *Usecase) Register(user auth.NewUser) error {
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

	err = uc.r.Register(user.Name, user.Email, hash)
	if err != nil {
		return exceptions.ServerError
	}

	return nil
}

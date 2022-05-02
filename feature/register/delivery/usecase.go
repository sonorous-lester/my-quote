package delivery

import (
	"myquote/domain"
	"myquote/domain/common"
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

}

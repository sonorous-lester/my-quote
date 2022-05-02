package delivery

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"myquote/domain/exceptions"
	"myquote/domain/register"
	"myquote/service/logger"
	"testing"
)

type MockedRegisterRepo struct {
	mock.Mock
}

func (m *MockedRegisterRepo) Find(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

type MockedEmailValidator struct {
	mock.Mock
}

func (m *MockedEmailValidator) Validate(s string) bool {
	args := m.Called(s)
	return args.Bool(0)
}

type MockedPasswordValidator struct {
	mock.Mock
}

func (m *MockedPasswordValidator) Validate(s string) bool {
	args := m.Called(s)
	return args.Bool(0)
}

type RegisterUsecaseTestSuite struct {
	suite.Suite
	uc   register.Usecase
	repo *MockedRegisterRepo
	pv   *MockedPasswordValidator
	ev   *MockedEmailValidator
}

func TestNewRegisterUsecase(t *testing.T) {
	suite.Run(t, new(RegisterUsecaseTestSuite))
}

func (s *RegisterUsecaseTestSuite) SetupTest() {
	l := logger.NewLogger("")
	s.repo = new(MockedRegisterRepo)
	s.pv = new(MockedPasswordValidator)
	s.ev = new(MockedEmailValidator)
	s.uc = NewRegisterUsecase(l, s.repo, s.pv, s.ev)
}

func (s *RegisterUsecaseTestSuite) TestRegisterInvalidEmailAddr() {
	user := register.NewUser{
		Email:    "123jkljl",
		Password: "dfadfjkladsf",
	}

	s.ev.On("Validate", user.Email).Return(false)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.InvalidEmailAddr, err)
}

func (s *RegisterUsecaseTestSuite) TestRegisterInvalidPasswordLength() {
	user := register.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklfdasfdsfadsfadsf",
	}
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(false)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.InvalidPasswordLength, err)
}

func (s *RegisterUsecaseTestSuite) TestRegisterUserExists() {
	user := register.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklf",
	}
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(true)
	s.repo.On("Find", user.Email).Return(true)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.UserExists, err)
}

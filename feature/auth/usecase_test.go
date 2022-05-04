package auth

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"myquote/domain/auth"
	"myquote/domain/exceptions"
	"myquote/domain/models"
	"myquote/service/logger"
	"testing"
)

type MockedAuthRepo struct {
	mock.Mock
}

func (m *MockedAuthRepo) Register(name string, email string, password string) error {
	args := m.Called(name, email, password)
	return args.Error(0)
}

func (m *MockedAuthRepo) FindUser(email string) (bool, models.UserModel, error) {
	args := m.Called(email)
	return args.Bool(0), args.Get(1).(models.UserModel), args.Error(2)
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

type MockedHashValidator struct {
	mock.Mock
}

func (m *MockedHashValidator) Hash(s string) (string, error) {
	args := m.Called(s)
	return args.String(0), args.Error(1)
}

func (m *MockedHashValidator) Compare(s string, h string) bool {
	args := m.Called(s, h)
	return args.Bool(0)
}

type MockedTokenGenerator struct {
	mock.Mock
}

func (m *MockedTokenGenerator) New() string {
	args := m.Called()
	return args.String(0)
}

type AuthUsecaseTestSuite struct {
	suite.Suite
	uc     auth.Usecase
	repo   *MockedAuthRepo
	pv     *MockedPasswordValidator
	ev     *MockedEmailValidator
	hashv  *MockedHashValidator
	tokeng *MockedTokenGenerator
}

func TestNewAuthUsecase(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}

func (s *AuthUsecaseTestSuite) SetupTest() {
	l := logger.NewLogger("")
	s.repo = new(MockedAuthRepo)
	s.pv = new(MockedPasswordValidator)
	s.ev = new(MockedEmailValidator)
	s.hashv = new(MockedHashValidator)
	s.tokeng = new(MockedTokenGenerator)
	s.uc = NewUsecase(l, s.repo, s.pv, s.ev, s.hashv, s.tokeng)
}

func (s *AuthUsecaseTestSuite) TestRegisterInvalidEmailAddr() {
	user := auth.NewUser{
		Email:    "123jkljl",
		Password: "dfadfjkladsf",
	}

	s.ev.On("Validate", user.Email).Return(false)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.InvalidEmailAddr, err)
}

func (s *AuthUsecaseTestSuite) TestRegisterInvalidPasswordLength() {
	user := auth.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklfdasfdsfadsfadsf",
	}
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(false)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.InvalidPasswordLength, err)
}

func (s *AuthUsecaseTestSuite) TestRegisterUserExists() {
	user := auth.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklf",
	}
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(true)
	s.repo.On("FindUser", user.Email).Return(true, models.UserModel{}, nil)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.UserExists, err)
}

func (s *AuthUsecaseTestSuite) TestRegisterThrowServerError() {
	user := auth.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklf",
	}
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(true)
	s.repo.On("FindUser", user.Email).Return(false, models.UserModel{}, exceptions.ServerError)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.ServerError, err)
}

func (s *AuthUsecaseTestSuite) TestRegisterThrowServerErrorHashPasswordFailure() {
	user := auth.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklf",
	}
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(true)
	s.repo.On("FindUser", user.Email).Return(false, models.UserModel{}, nil)
	s.hashv.On("Hash", user.Password).Return("", exceptions.ServerError)
	s.repo.On("Register", user.Name, user.Email, user.Password).Return(nil)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.ServerError, err)
}

func (s *AuthUsecaseTestSuite) TestRegisterThrowServerErrorWhenRegisterFailure() {
	user := auth.NewUser{
		Email:    "123@gmail.com",
		Password: "dfadfjklf",
	}
	hash := "hashresult"
	s.ev.On("Validate", user.Email).Return(true)
	s.pv.On("Validate", user.Password).Return(true)
	s.repo.On("FindUser", user.Email).Return(false, models.UserModel{}, nil)
	s.hashv.On("Hash", user.Password).Return(hash, nil)
	s.repo.On("Register", user.Name, user.Email, hash).Return(exceptions.ServerError)
	err := s.uc.Register(user)

	s.Assert().Equal(exceptions.ServerError, err)
}

func (s *AuthUsecaseTestSuite) TestLoginThrowUserNotExistException() {
	info := auth.LoginInfo{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	s.repo.On("FindUser", info.Email).Return(false, models.UserModel{}, exceptions.UserNotExists)
	_, err := s.uc.Login(info)
	s.Assert().Equal(exceptions.UserNotExists, err)
}

func (s *AuthUsecaseTestSuite) TestLoginThrowUserServerErrorException() {
	info := auth.LoginInfo{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	s.repo.On("FindUser", info.Email).Return(true, models.UserModel{}, exceptions.ServerError)
	_, err := s.uc.Login(info)
	s.Assert().Equal(exceptions.ServerError, err)
}

func (s *AuthUsecaseTestSuite) TestLoginThrowAuthErrorException() {
	info := auth.LoginInfo{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	user := models.UserModel{Password: "this is a hash"}

	s.repo.On("FindUser", info.Email).Return(true, user, nil)
	s.hashv.On("Compare", info.Password, user.Password).Return(false)
	_, err := s.uc.Login(info)
	s.Assert().Equal(exceptions.AuthError, err)
}

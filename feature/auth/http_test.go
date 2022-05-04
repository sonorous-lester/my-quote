package auth

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"myquote/domain"
	"myquote/domain/auth"
	"myquote/domain/common"
	"myquote/domain/exceptions"
	"myquote/domain/models"
	"myquote/service/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockedAuthUsecase struct {
	mock.Mock
}

func (m *MockedAuthUsecase) Signout() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedAuthUsecase) Register(user auth.NewUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockedAuthUsecase) Login(i auth.Anonymous) (models.User, error) {
	args := m.Called(i)
	return args.Get(0).(models.User), args.Error(1)
}

type AuthTestSuite struct {
	suite.Suite
	uc *MockedAuthUsecase
	l  domain.Logger
	g  *gin.Engine
	r  *httptest.ResponseRecorder
}

func TestAuthHTTPHandler(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupTest() {
	s.uc = new(MockedAuthUsecase)
	s.l = logger.NewLogger("")
	s.g = gin.Default()
	s.r = httptest.NewRecorder()
}

func newTestRequest(method string, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

func (s *AuthTestSuite) TestRegisterSuccess() {
	newUser := auth.NewUser{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	body, _ := json.Marshal(newUser)
	s.uc.On("Register", newUser).Return(nil)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, REGISTER_ENDPOINT, body)
	s.g.ServeHTTP(s.r, req)
	s.Assert().Equal(http.StatusOK, s.r.Code)
}

func (s *AuthTestSuite) TestRegisterInvalidInput() {
	s.uc.On("Register", nil).Return(nil)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, REGISTER_ENDPOINT, nil)
	s.g.ServeHTTP(s.r, req)

	s.Assert().Equal(http.StatusBadRequest, s.r.Code)
	var m common.Message
	json.Unmarshal(s.r.Body.Bytes(), &m)
	s.Assert().Equal(exceptions.InvalidInput.Error(), m.Message)
}

func (s *AuthTestSuite) TestRegisterShowInvalidMessage() {
	newUser := auth.NewUser{
		Email:    "123@",
		Password: "123456",
	}
	body, _ := json.Marshal(newUser)
	s.uc.On("Register", newUser).Return(exceptions.InvalidEmailAddr)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, REGISTER_ENDPOINT, body)
	s.g.ServeHTTP(s.r, req)

	s.Assert().Equal(http.StatusBadRequest, s.r.Code)
	var m common.Message
	json.Unmarshal(s.r.Body.Bytes(), &m)
	s.Assert().Equal(exceptions.InvalidEmailAddr.Error(), m.Message)
}

func (s *AuthTestSuite) TestLoginSuccess() {
	info := auth.Anonymous{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	body, _ := json.Marshal(info)

	user := models.User{
		ID:        1,
		Name:      "Lester",
		Email:     "123@gmail.com",
		Token:     "secret token",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	s.uc.On("Login", info).Return(user, nil)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, LOGIN_ENDPOINT, body)
	s.g.ServeHTTP(s.r, req)

	var u models.User
	json.Unmarshal(s.r.Body.Bytes(), &u)
	s.Assert().Equal(http.StatusOK, s.r.Code)
	s.Assert().Equal(user.Name, u.Name)
	s.Assert().Equal(user.Email, u.Email)
	s.Assert().Equal(user.Token, u.Token)
}

func (s *AuthTestSuite) TestShowLoginInvalidInputException() {
	s.uc.On("Login", mock.Anything).Return(models.User{}, nil)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, LOGIN_ENDPOINT, nil)
	s.g.ServeHTTP(s.r, req)

	var m common.Message
	json.Unmarshal(s.r.Body.Bytes(), &m)
	s.Assert().Equal(http.StatusBadRequest, s.r.Code)
	s.Assert().Equal(exceptions.InvalidInput.Error(), m.Message)
}

func (s *AuthTestSuite) TestLoginShowBadRequestException() {
	info := auth.Anonymous{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	body, _ := json.Marshal(info)

	s.uc.On("Login", info).Return(models.User{}, exceptions.AuthError)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, LOGIN_ENDPOINT, body)
	s.g.ServeHTTP(s.r, req)

	var m common.Message
	json.Unmarshal(s.r.Body.Bytes(), &m)
	s.Assert().Equal(http.StatusBadRequest, s.r.Code)
	s.Assert().Equal(exceptions.AuthError.Error(), m.Message)
}

func (s *AuthTestSuite) TestLoginShowServerErrorException() {
	info := auth.Anonymous{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	body, _ := json.Marshal(info)

	s.uc.On("Login", info).Return(models.User{}, exceptions.ServerError)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, LOGIN_ENDPOINT, body)
	s.g.ServeHTTP(s.r, req)

	var m common.Message
	json.Unmarshal(s.r.Body.Bytes(), &m)
	s.Assert().Equal(http.StatusInternalServerError, s.r.Code)
	s.Assert().Equal(exceptions.ServerError.Error(), m.Message)
}

func (s *AuthTestSuite) TestSignoutSuccess() {
	s.uc.On("Signout").Return(nil)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, SIGNOUT_ENDPOINT, nil)
	s.g.ServeHTTP(s.r, req)
	s.Assert().Equal(http.StatusOK, s.r.Code)
}

func (s *AuthTestSuite) TestRespondServerErrorWhenSignoutFailure() {
	s.uc.On("Signout").Return(exceptions.ServerError)
	NewAuthHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, SIGNOUT_ENDPOINT, nil)
	s.g.ServeHTTP(s.r, req)
	s.Assert().Equal(http.StatusInternalServerError, s.r.Code)
}

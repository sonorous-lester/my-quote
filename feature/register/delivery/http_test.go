package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"myquote/domain"
	"myquote/domain/common"
	"myquote/domain/exceptions"
	"myquote/domain/register"
	"myquote/service/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockedRegisterUsecase struct {
	mock.Mock
}

func (m *MockedRegisterUsecase) Register(user register.NewUser) error {
	args := m.Called(user)
	return args.Error(0)
}

type RegisterTestSuite struct {
	suite.Suite
	uc *MockedRegisterUsecase
	l  domain.Logger
	g  *gin.Engine
	r  *httptest.ResponseRecorder
}

func TestRegisterHTTPHandler(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) SetupTest() {
	s.uc = new(MockedRegisterUsecase)
	s.l = logger.NewLogger("")
	s.g = gin.Default()
	s.r = httptest.NewRecorder()
}

func newTestRequest(method string, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

func (s *RegisterTestSuite) TestRegisterSuccess() {
	newUser := register.NewUser{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	body, _ := json.Marshal(newUser)
	s.uc.On("Register", newUser).Return(nil)
	NewRegisterHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, REGISTER_ENDPOINT, body)
	s.g.ServeHTTP(s.r, req)
	s.Assert().Equal(http.StatusOK, s.r.Code)
}

func (s *RegisterTestSuite) TestRegisterInvalidInput() {
	s.uc.On("Register", nil).Return(nil)
	NewRegisterHTTPHandler(s.g, s.l, s.uc)
	req, _ := newTestRequest(http.MethodPost, REGISTER_ENDPOINT, nil)
	s.g.ServeHTTP(s.r, req)

	s.Assert().Equal(http.StatusBadRequest, s.r.Code)
	var m common.Message
	json.Unmarshal(s.r.Body.Bytes(), &m)
	s.Assert().Equal(exceptions.ValidInput.Error(), m.Message)
}

package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

var l = logger.NewLogger("")

func TestRegisterSuccess(t *testing.T) {
	newUser := register.NewUser{
		Email:    "123@gmail.com",
		Password: "123456",
	}
	body, _ := json.Marshal(newUser)

	g := gin.Default()
	mockedUc := new(MockedRegisterUsecase)
	mockedUc.On("Register", newUser).Return(nil)
	NewRegisterHTTPHandler(g, l, mockedUc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	g.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

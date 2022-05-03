package auth

import (
	"github.com/gin-gonic/gin"
	"myquote/domain"
	"myquote/domain/auth"
	"myquote/domain/common"
	"myquote/domain/exceptions"
	"net/http"
)

type handler struct {
	logger     domain.Logger
	registerUc auth.Usecase
}

const REGISTER_ENDPOINT = "/api/auth"
const LOGIN_ENDPOINT = "/api/login"

func NewAuthHTTPHandler(c *gin.Engine, l domain.Logger, uc auth.Usecase) {
	handler := &handler{logger: l, registerUc: uc}
	c.POST(REGISTER_ENDPOINT, handler.register)
	c.POST(LOGIN_ENDPOINT, handler.login)
}

func (h *handler) register(c *gin.Context) {
	var user auth.NewUser
	err := c.Bind(&user)
	if err != nil {
		h.logger.Debugf("Convert new user json error: %s", err.Error())
		c.JSON(http.StatusBadRequest, common.Message{Message: exceptions.InvalidInput.Error()})
		return
	}
	err = h.registerUc.Register(user)
	if err != nil {
		h.logger.Warnf("auth user error: %s", err.Error())
		c.JSON(http.StatusBadRequest, common.Message{Message: err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (h *handler) login(c *gin.Context) {
	var info auth.LoginInfo
	c.Bind(&info)
	user, _ := h.registerUc.Login(info)
	c.JSON(http.StatusOK, user)
}

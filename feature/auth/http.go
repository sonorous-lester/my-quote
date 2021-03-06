package auth

import (
	"errors"
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
const SIGNOUT_ENDPOINT = "/api/signout"

func NewAuthHTTPHandler(c *gin.Engine, l domain.Logger, uc auth.Usecase) {
	handler := &handler{logger: l, registerUc: uc}
	c.POST(REGISTER_ENDPOINT, handler.register)
	c.POST(LOGIN_ENDPOINT, handler.login)
	c.POST(SIGNOUT_ENDPOINT, handler.signout)
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
	var info auth.Anonymous
	err := c.Bind(&info)
	if err != nil {
		h.logger.Debugf("Convert login info json error: %s", err.Error())
		c.JSON(http.StatusBadRequest, common.Message{Message: exceptions.InvalidInput.Error()})
		return
	}
	user, err := h.registerUc.Login(info)
	if err != nil && errors.Is(err, exceptions.ServerError) {
		c.JSON(http.StatusInternalServerError, common.Message{Message: err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, common.Message{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) signout(c *gin.Context) {
	err := h.registerUc.Signout()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Message{Message: "sign out successful"})
}

package delivery

import (
	"github.com/gin-gonic/gin"
	"myquote/domain"
	"myquote/domain/common"
	"myquote/domain/exceptions"
	"myquote/domain/register"
	"net/http"
)

type handler struct {
	logger     domain.Logger
	registerUc register.Usecase
}

const REGISTER_ENDPOINT = "/api/register"

func NewRegisterHTTPHandler(c *gin.Engine, l domain.Logger, uc register.Usecase) {
	handler := &handler{logger: l, registerUc: uc}
	c.POST(REGISTER_ENDPOINT, handler.register)
}

func (h *handler) register(c *gin.Context) {
	var user register.NewUser
	err := c.Bind(&user)
	if err != nil {
		h.logger.Debugf("Convert new user json error: %s", err.Error())
		c.JSON(http.StatusBadRequest, common.Message{Message: exceptions.InvalidInput.Error()})
		return
	}
	err = h.registerUc.Register(user)
	if err != nil {
		h.logger.Warnf("register user error: %s", err.Error())
		c.JSON(http.StatusBadRequest, common.Message{Message: err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

package delivery

import (
	"github.com/gin-gonic/gin"
	"myquote/domain"
	"myquote/domain/register"
	"net/http"
)

type handler struct {
	logger     domain.Logger
	registerUc register.Usecase
}

func NewRegisterHTTPHandler(c *gin.Engine, l domain.Logger, uc register.Usecase) {
	handler := &handler{logger: l, registerUc: uc}
	c.POST("/api/register", handler.register)
}

func (h *handler) register(c *gin.Context) {
	var user register.NewUser
	err := c.Bind(&user)
	if err != nil {
		h.logger.Debugf("Convert new user json error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	h.registerUc.Register(user)
	c.Writer.WriteHeader(http.StatusOK)
}

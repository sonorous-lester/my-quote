package delivery

import (
	"github.com/gin-gonic/gin"
	"myquote/domain"
	"myquote/domain/register"
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

}


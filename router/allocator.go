package router

import (
	"numbers-allocator/control"
	"numbers-allocator/model"

	"github.com/gin-gonic/gin"
	"net/http"
)

func initAllocatorRouter(engine *gin.Engine) {

	engine.POST(model.ApiApplyIdRange, func(c *gin.Context) {
		c.JSON(http.StatusOK, control.Allocator.ApplyIdRange(c))
	})
}

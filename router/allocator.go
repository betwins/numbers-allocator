package router

import (
	"context"
	"numbers-allocator/control"
	"numbers-allocator/model"
	"numbers-allocator/trx"

	"github.com/gin-gonic/gin"
	"net/http"
)

func initAllocatorRouter(engine *gin.Engine) {

	engine.POST(model.ApiApplyIdRange, func(c *gin.Context) {
		c.JSON(http.StatusOK, trx.Transaction(context.Background(), c, control.Allocator.ApplyIdRange))
	})
}

package control

import (
	"numbers-allocator/errcode"
	"numbers-allocator/model"
	"numbers-allocator/service"

	"github.com/gin-gonic/gin"
	"github.com/maczh/mgin/i18n"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/models"
	"time"
)

type allocatorController struct{}

var Allocator allocatorController

// ApplyIdRange	godoc
// @Summary		独占号段申请
// @Description	独占号段申请
// @Tags	号段管理
// @Accept	application/json
// @Produce json
// @Param X-Lang header string false "语言"
// @Param Partner-Id header string true "代理商id"
// @Param params body model.ApplyReq true "请求体"
// @Success 200 {object} model.NewRangeResp
// @Router	/ids/apply [post]
func (c *allocatorController) ApplyIdRange(ctx *gin.Context) models.Result[any] {
	var req model.ApplyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logs.Error("参数绑定错误 {}", err.Error())
		return models.Error(-1, err.Error())
	}

	if req.Type == 0 {
		return i18n.ParamError("type")
	}
	if req.Step == 0 {
		return i18n.ParamError("step")
	}
	if req.Day == "" {
		return i18n.ParamError("step")
	}

	applyDay, err := time.Parse("2006-01-02", req.Day)
	if err != nil {
		return errcode.ParamFormatErr.MGErrorWithArgs("day")
	}

	rangeStart, rangeEnd, err := service.Allocator.GetIdRange(applyDay, req.Type, req.Step)
	if err != nil {
		return models.Error(-1, err.Error())
	}

	resp := model.NewRangeResp{
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
	}
	return i18n.Success[any](resp)

}

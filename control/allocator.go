package control

import (
	"context"
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
// @Router	/numbers/apply [post]
func (c *allocatorController) ApplyIdRange(ctx context.Context, ginCtx *gin.Context) models.Result[any] {
	var req model.ApplyReq
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		logs.Error("参数绑定错误 {}", err.Error())
		return models.Error(-1, err.Error())
	}

	if req.AppName == "" {
		return i18n.ParamError("appName")
	}
	if len(req.AppName) > 50 {
		return errcode.AppNameTooLong.MGError()
	}
	if req.BizType == "" {
		return i18n.ParamError("bizType")
	}
	if len(req.BizType) > 50 {
		return errcode.BizTypeTooLong.MGError()
	}
	if req.Step <= 0 {
		return i18n.ParamError("step")
	}
	if req.Day == "" {
		return i18n.ParamError("day")
	}
	if len(req.Day) != 8 {
		return errcode.DayFormatErr.MGErrorWithArgs("YYYYMMDD")
	}

	applyDay, err := time.Parse("20060102", req.Day)
	if err != nil {
		return errcode.DayFormatErr.MGErrorWithArgs("YYYYMMDD")
	}

	logs.Debug("申请号段 applyDay: {} req: {}", applyDay, req)

	rangeStart, rangeEnd, err := service.Allocator.GetIdRange(ctx, &req)
	if err != nil {
		return models.Error(-1, err.Error())
	}

	resp := model.NewRangeResp{
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
	}
	return i18n.Success[any](resp)

}

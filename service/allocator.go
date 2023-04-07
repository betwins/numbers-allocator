package service

import (
	"context"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/middleware/trace"
	"numbers-allocator/dao"
	"numbers-allocator/errcode"
	"numbers-allocator/model"
)

type allocatorService struct{}

var Allocator allocatorService

func (s *allocatorService) GetIdRange(ctx context.Context, req *model.ApplyReq) (int64, int64, error) {

	var rangeStart int64
	var rangeEnd int64

	exist, err := dao.Allocator.GetEntity(ctx, req.AppName, req.BizType, req.Day)
	if err != nil {
		logs.Error("{} {} {} 查询失败 err: {}, req: {}", trace.GetGoroutineID(), req.AppName, req.BizType, err.Error(), req)
		return 0, 0, errcode.DbQueryErr.Error()
	}

	if exist == nil {
		//初始化号段记录
		rangeStart, rangeEnd = getDayInitRange(req.Step)
		newEntity := model.Allocator{
			ApplyAppName:   req.AppName,
			ApplyBizType:   req.BizType,
			CurrentStartId: rangeStart,
			IncrementStep:  req.Step,
			ApplyDate:      req.Day,
			Version:        1,
		}

		bNewInsert, err := dao.Allocator.AddEntity(ctx, &newEntity)
		if err != nil {
			//出错
			logs.Error("{} {} {} 创建号段记录失败 err: {}, newEntity: {} req: {}", trace.GetGoroutineID(), req.AppName, req.BizType, err.Error(), newEntity, req)
			return 0, 0, err
		}
		if bNewInsert {
			logs.Debug("{} {} {} 新创建当天号段记录 {}", trace.GetGoroutineID(), req.AppName, req.BizType, newEntity)
			return rangeStart, rangeEnd, nil
		}
	}

	mustExist, err := dao.Allocator.GetEntity(ctx, req.AppName, req.BizType, req.Day)
	if err != nil {
		logs.Error("{} {} {} 查询失败 err: {}, req: {}", trace.GetGoroutineID(), req.AppName, req.BizType, err.Error(), req)
		return 0, 0, errcode.DbQueryErr.Error()
	}
	if mustExist == nil {
		logs.Error("{} {} {} 不应该不存在记录 req: {}", trace.GetGoroutineID(), req.AppName, req.BizType, req)
		return 0, 0, errcode.DbQueryErr.Error()
	}

	logs.Debug("{} {} {} 当天非第一次申请， 原号段: {} 新请求: {}", trace.GetGoroutineID(), req.AppName, req.BizType, mustExist, req)

	//更新对应申请日期使用的号段
	err = dao.Allocator.UpdateRange(ctx, req.Step, mustExist.Id)
	if err != nil {
		//出错
		logs.Error("更新记录失败 {}", err.Error())
		return 0, 0, err
	}

	entityAfterUpdated, err := dao.Allocator.GetEntityById(ctx, mustExist.Id)
	if err != nil {
		//出错
		logs.Error("{} {} {} 查询更新后的号段失败 err: {} 请求: {}", trace.GetGoroutineID(), req.AppName, req.BizType, err.Error(), req)
		return 0, 0, err
	}

	newRangeStart := entityAfterUpdated.CurrentStartId
	newRangeEnd := entityAfterUpdated.CurrentStartId + int64(entityAfterUpdated.IncrementStep)

	logs.Debug("{} {} {} 当天非第一次申请结果，返回新号段:{} {}, 更新后记录:{}", trace.GetGoroutineID(), req.AppName, req.BizType, newRangeStart, newRangeEnd, entityAfterUpdated)

	return newRangeStart, newRangeEnd, nil
}

func getDayInitRange(step int) (int64, int64) {
	rangeStart := 32619 //先采用固定值验证
	//rangeStart := rand.Intn(49873) + 126
	rangeEnd := rangeStart + step
	return int64(rangeStart), int64(rangeEnd)
}

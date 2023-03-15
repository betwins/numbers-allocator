package service

import (
	"github.com/maczh/mgin/logs"
	"numbers-allocator/dao"
	"numbers-allocator/errcode"
	"numbers-allocator/model"
)

type allocatorService struct{}

var Allocator allocatorService

func (s *allocatorService) GetIdRange(req *model.ApplyReq) (int64, int64, error) {

	var rangeStart int64
	var rangeEnd int64

	entity, err := dao.Allocator.GetEntity(req.AppName, req.BizType, req.Day)
	if err != nil {
		logs.Error("{} {} 查询失败 err: {}, req: {}", req.AppName, req.BizType, err.Error(), req)
		return 0, 0, errcode.DbQueryErr.Error()
	}
	if entity == nil {
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

		err = dao.Allocator.AddEntity(&newEntity)
		if err != nil {
			//出错
			logs.Error("{} {} 创建号段记录失败 err: {}, newEntity: {} req: {}", req.AppName, req.BizType, err.Error(), newEntity, req)
			return 0, 0, err
		}
		entityAfterUpdated, err := dao.Allocator.GetEntity(req.AppName, req.BizType, req.Day)
		if err != nil {
			//出错
			logs.Error("{} {} 查询更新后的号段失败 err: {} 请求: {}, 进新更新的数据: {}", req.AppName, req.BizType, err.Error(), req, newEntity)
			return 0, 0, err
		}
		logs.Debug("{} {} 当天第一个号段 请求:{} 返回号段: {} {}， 更新后号段记录: {}", req.AppName, req.BizType, req, rangeStart, rangeEnd, entityAfterUpdated)
		return rangeStart, rangeEnd, nil
	}

	logs.Debug("{} {} 当天非第一次申请，原号段记录: {}, 新请求: {}", req.AppName, req.BizType, entity, req)
	//更新对应申请日期使用的号段
	rangeStart = entity.CurrentStartId + int64(entity.IncrementStep) + 1
	rangeEnd = rangeStart + int64(req.Step)
	entity.CurrentStartId = rangeStart
	entity.IncrementStep = req.Step
	oldVersion := entity.Version
	entity.Version = entity.Version + 1
	err = dao.Allocator.UpdateEntity(entity, oldVersion)
	if err != nil {
		//出错
		logs.Error("更新记录失败")
		return 0, 0, err
	}

	entityAfterUpdated, err := dao.Allocator.GetEntity(req.AppName, req.BizType, req.Day)
	if err != nil {
		//出错
		logs.Error("{} {} 查询更新后的号段失败 err: {} 请求: {}, 更新数据: {}", req.AppName, req.BizType, err.Error(), req, entity)
		return 0, 0, err
	}

	logs.Debug("{} {} 当天非第一次申请结果，返回新号段:{} {}, 更新数据:{}, 更新后记录:{}", req.AppName, req.BizType, rangeStart, rangeEnd, entity, entityAfterUpdated)

	return rangeStart, rangeEnd, nil
}

func getDayInitRange(step int) (int64, int64) {
	rangeStart := 32619 //先采用固定值验证
	//rangeStart := rand.Intn(49873) + 126
	rangeEnd := rangeStart + step
	return int64(rangeStart), int64(rangeEnd)
}

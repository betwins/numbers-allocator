package service

import (
	"github.com/maczh/mgin/db"
	"github.com/maczh/mgin/logs"
	"math/rand"
	"numbers-allocator/dao"
	"numbers-allocator/errcode"
	"numbers-allocator/model"
	"numbers-allocator/multidb"
)

type allocatorService struct{}

var Allocator allocatorService

func (s *allocatorService) GetIdRange(req *model.ApplyReq) (int64, int64, error) {

	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		logs.Error("获取mysql连接失败 err: {}", err.Error())
		return 0, 0, errcode.DbConnectErr.Error()
	}

	var rangeStart int64
	var rangeEnd int64

	entity, err := dao.Allocator.GetEntity(req.AppName, req.BizType, req.Day)
	if err != nil {
		logs.Error("查询失败 err: {}, req: {}", err.Error(), req)
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
		err = conn.Debug().Create(&newEntity).Error
		if err != nil {
			//出错
			logs.Error("创建号段记录失败 err: {}, newEntity: {} req: {}", err.Error(), newEntity, req)
			return 0, 0, err
		}
		return rangeStart, rangeEnd, nil
	}

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

	logs.Debug("获取到了新号段 请求：{} {} {} {}", req.AppName, req.BizType, req.Day, req.Step)
	logs.Debug("获取到了新号段 结果：{} {} {} {} {} {}", entity.CurrentStartId, entity.IncrementStep, entity.Version, entity.ApplyDate, rangeStart, rangeEnd)
	return rangeStart, rangeEnd, nil
}

func getDayInitRange(step int) (int64, int64) {
	rangeStart := rand.Intn(49873) + 126
	rangeEnd := rangeStart + step
	return int64(rangeStart), int64(rangeEnd)
}

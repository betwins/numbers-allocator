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

func (s *allocatorService) GetIdRange(req *model.ApplyReq) (int, int, error) {

	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		return 0, 0, errcode.DbConnectErr.Error()
	}

	var rangeStart int
	var rangeEnd int

	entity, err := dao.Allocator.GetEntity(req.AppName, req.BizType, req.Day)
	if err != nil {
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
			logs.Error("创建号段记录失败 {}, {}", err.Error(), newEntity)
			return 0, 0, err
		}
		return rangeStart, rangeEnd, nil
	}

	//更新对应申请日期使用的号段
	rangeStart = entity.CurrentStartId + entity.IncrementStep + 1
	rangeEnd = rangeStart + req.Step
	entity.CurrentStartId = rangeStart
	entity.IncrementStep = req.Step
	entity.ApplyDate = req.Day
	oldVersion := entity.Version
	entity.Version = entity.Version + 1
	err = dao.Allocator.UpdateEntity(entity, oldVersion)
	if err != nil {
		//出错
		return 0, 0, err
	}
	return rangeStart, rangeEnd, nil
}

func getDayInitRange(step int) (int, int) {
	//todo 先用从固定值开始的号段进行验证，后面改回随机数开始号段
	//rangeStart := 32178
	//rangeEnd := rangeStart + step
	rangeStart := rand.Intn(49873) + 126
	rangeEnd := rangeStart + step
	return rangeStart, rangeEnd
}

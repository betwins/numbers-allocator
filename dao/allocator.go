package dao

import (
	"errors"
	"fmt"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/middleware/trace"
	"numbers-allocator/errcode"
	"numbers-allocator/model"
	"numbers-allocator/multidb"

	"github.com/maczh/mgin/db"
	"gorm.io/gorm"
)

type allocatorDao struct{}

var Allocator allocatorDao

func (s *allocatorDao) AddEntity(entity *model.Allocator) error {

	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		return errcode.DbConnectErr.Error()
	}

	err = conn.Debug().Model(entity).Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *allocatorDao) UpdateEntity(entity *model.Allocator, oldVersion int) error {
	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		return errcode.DbConnectErr.Error()
	}

	logs.Debug("获取新号段，更新新号段 {} {} {} {}", entity.Id, oldVersion, entity.Version, trace.GetGoroutineID())
	//err = conn.Debug().Model(entity).Where("id = ? and version = ?", entity.Id, oldVersion).Updates(entity).Error

	sql := fmt.Sprintf("update numbers_allocator set current_start_id = #v and set increment_step = #v and set version = #v where id = #v and version = #v",
		entity.CurrentStartId, entity.IncrementStep, entity.Version, entity.Id, oldVersion)
	err = conn.Exec(sql).Error
	if err != nil {
		logs.Error("获取新号段，更新失败 {}", trace.GetGoroutineID())
		return err
	}
	logs.Debug("获取新号段, 更新成功 {}", trace.GetGoroutineID())

	return nil
}

func (s *allocatorDao) GetEntity(appName string, bizType string, day string) (*model.Allocator, error) {
	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		return nil, errcode.DbConnectErr.Error()
	}

	var entity model.Allocator

	err = conn.Debug().Model(&entity).Where("apply_app_name = ? and apply_biz_type = ? and apply_date =?", appName, bizType, day).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

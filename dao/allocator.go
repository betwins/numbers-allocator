package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/middleware/trace"
	"numbers-allocator/errcode"
	"numbers-allocator/model"
	"numbers-allocator/trx"

	"gorm.io/gorm"
)

type allocatorDao struct{}

var Allocator allocatorDao

func (s *allocatorDao) AddEntity(ctx context.Context, entity *model.Allocator) (bool, error) {

	conn, err := trx.ExtractDb(ctx)
	if err != nil {
		return false, errcode.DbConnectErr.Error()
	}

	sqlFormat := "INSERT IGNORE INTO numbers_allocator (`apply_app_name`, `apply_biz_type`, `current_start_id`, `increment_step`, `apply_date`, `version`) VALUES('%s', '%s', %d, %d, '%s', %d)"
	sql := fmt.Sprintf(sqlFormat, entity.ApplyAppName, entity.ApplyBizType, entity.CurrentStartId, entity.IncrementStep, entity.ApplyDate, entity.Version)

	result := conn.Debug().Exec(sql)
	if result.Error != nil {
		logs.Error("{} 获取新号段，更新失败", trace.GetGoroutineID())
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		//已存在
		return false, nil
	}

	return true, nil
}

func (s *allocatorDao) UpdateRange(ctx context.Context, step int, id int) error {
	conn, err := trx.ExtractDb(ctx)
	if err != nil {
		return errcode.DbConnectErr.Error()
	}

	//logs.Debug("{} 获取新号段，更新新号段 {} {} {} {}", trace.GetGoroutineID(), entity.Id, oldVersion, entity.Version)
	//result := conn.Debug().Model(entity).Where("id = ? and version = ?", entity.Id, oldVersion).Updates(entity)

	sql := fmt.Sprintf("update numbers_allocator set current_start_id = current_start_id + increment_step + 1, increment_step = %d, version=version+1 where id = %d", step, id)

	result := conn.Debug().Exec(sql)
	if result.Error != nil {
		logs.Error("{} 获取新号段，更新失败", trace.GetGoroutineID())
		return result.Error
	}
	if result.RowsAffected == 0 {
		logs.Error("{} 获取号段并发冲突 {} {}", trace.GetGoroutineID(), step, id)
		return errcode.ConcurrencyConflict.Error()
	}
	logs.Debug("{} 获取新号段, 更新成功 {}", trace.GetGoroutineID(), step)

	return nil
}

func (s *allocatorDao) UpdateEntity(ctx context.Context, entity *model.Allocator, oldVersion int) error {
	conn, err := trx.ExtractDb(ctx)
	if err != nil {
		return errcode.DbConnectErr.Error()
	}

	logs.Debug("{} 获取新号段，更新新号段 {} {} {} {}", trace.GetGoroutineID(), entity.Id, oldVersion, entity.Version)
	result := conn.Debug().Model(entity).Where("id = ? and version = ?", entity.Id, oldVersion).Updates(entity)

	if result.Error != nil {
		logs.Error("{} 获取新号段，更新失败", trace.GetGoroutineID())
		return result.Error
	}
	if result.RowsAffected == 0 {
		logs.Error("{} 获取号段并发冲突 {} {}", trace.GetGoroutineID(), entity, oldVersion)
		return errcode.ConcurrencyConflict.Error()
	}
	logs.Debug("{} 获取新号段, 更新成功 {}", trace.GetGoroutineID(), entity)

	return nil
}

func (s *allocatorDao) GetEntity(ctx context.Context, appName string, bizType string, day string) (*model.Allocator, error) {

	conn, err := trx.ExtractDb(ctx)
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

func (s *allocatorDao) GetEntityById(ctx context.Context, id int) (*model.Allocator, error) {
	conn, err := trx.ExtractDb(ctx)
	if err != nil {
		return nil, errcode.DbConnectErr.Error()
	}

	var entity model.Allocator

	err = conn.Debug().Model(&entity).Where("id = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

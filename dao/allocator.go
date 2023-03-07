package dao

import (
	"errors"
	"numbers-allocator/errcode"
	"numbers-allocator/model"
	"numbers-allocator/multidb"

	"github.com/maczh/mgin/db"
	"gorm.io/gorm"
	"time"
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

	err = conn.Debug().Model(entity).Where("id = ? and version = ?", entity.Id, oldVersion).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *allocatorDao) GetEntity(applyType int, applyDay time.Time) (*model.Allocator, error) {
	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		return nil, errcode.DbConnectErr.Error()
	}

	var entity model.Allocator

	err = conn.Debug().Model(&entity).Where("apply_type = ? and apply_date =?", applyType, applyDay.Format("2006-01-02")).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

package model

import "time"

type Allocator struct {
	Id             int       `json:"id"  gorm:"column:id;primary_key;auto_increment"`
	ApplyType      int       `json:"applyType" gorm:"column:apply_type"`
	CurrentStartId int       `json:"currentStartId"  gorm:"column:current_start_id"`
	IncrementStep  int       `json:"incrementStep"  gorm:"column:increment_step"`
	ApplyDate      time.Time `json:"applyDate"  gorm:"column:apply_date"`
	Version        int       `json:"version" gorm:"column:version"`
}

func (*Allocator) TableName() string {
	return "id_range_allocator"
}

package model

type Allocator struct {
	Id             int    `json:"id"  gorm:"column:id;primary_key;auto_increment"`
	ApplyAppName   string `json:"applyAppName" gorm:"column:apply_app_name"`      //申请独占号段的应用名
	ApplyBizType   string `json:"applyBizType" gorm:"column:apply_biz_type"`      //申请独占号段的业务类型, 业务方需要确保appName + bizType 不与其它申请者重复
	CurrentStartId int64  `json:"currentStartId"  gorm:"column:current_start_id"` //当前申请号段起始值（最近一次分配的最大号段，已被占用）
	IncrementStep  int    `json:"incrementStep"  gorm:"column:increment_step"`    //当前号段步长
	ApplyDate      string `json:"applyDate"  gorm:"column:apply_date"`            //号段应用日期
	Version        int    `json:"version" gorm:"column:version"`                  //版本号
}

func (*Allocator) TableName() string {
	return "numbers_allocator"
}

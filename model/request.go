package model

type ApplyReq struct {
	AppName string `json:"appName"` //"申请应用名"
	BizType string `json:"bizType"` //应用内使用号段的业务类型，业务方需要确保appName + bizType 不与其它申请者重复
	Day     string `json:"day"`     //"日期格式: 20060102" 号段应用日期，获得的号段会确保该日期内独占（在appName+bizType范围内独点）
	Step    int    `json:"step"`    //"号段步长" 申请号段的步长, 建议申请步长为1000，或不超过100000
}

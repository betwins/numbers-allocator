package model

type ApplyReq struct {
	Type int    `json:"applyType"` //"申请类型"
	Day  string `json:"applyDay"`  //"日期格式: 20060102"
	Step int    `json:"step"`      //"号段步长"
}

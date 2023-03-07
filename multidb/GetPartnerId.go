package multidb

import (
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/middleware/trace"
)

func GetPartnerId() string {
	partnerId := trace.GetHeader("Partner-Id")
	if partnerId == "" {
		logs.Error("Partner-Id header缺失，多库错误")
	}
	return partnerId
}

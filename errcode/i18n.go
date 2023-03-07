package errcode

import (
	"errors"
	"fmt"
	"github.com/maczh/mgin/errcode"
	"github.com/maczh/mgin/i18n"
	"github.com/maczh/mgin/models"
)

type MessageId string

const (
	UrlNotFound        MessageId = "404"
	SystemError        MessageId = "系统异常"
	DbConnectErr       MessageId = "数据库连接失败"
	DbInsertErr        MessageId = "数据库插入失败"
	DbUpdateErr        MessageId = "数据库更新失败"
	DbDeleteErr        MessageId = "数据库删除失败"
	DataNotFound       MessageId = "数据库查无数据"
	ParamLost          MessageId = "参数不可为空"
	ParamError         MessageId = "参数错误"
	ParamFormatErr     MessageId = "参数{}格式或类型错误"
	ConnectFail        MessageId = "网络连接失败"
	ServiceUnavailable MessageId = "服务不存在"
	Success            MessageId = "success"
	DbQueryErr         MessageId = "数据库查询失败"
)

func (receiver MessageId) Error() error {
	return errors.New(i18n.String(string(receiver)))
}

func (receiver MessageId) ErrorWithMsg(msg string) error {
	return errors.New(fmt.Sprintf("%s: %s", i18n.String(string(receiver)), msg))
}

func (receiver MessageId) ErrorWithArgs(args ...any) error {
	return errors.New(i18n.Format(string(receiver), args...))
}

func (receiver MessageId) ErrorJoinMsgId(msgId MessageId) error {
	return receiver.ErrorWithMsg(msgId.GetText())
}

func (receiver MessageId) MGError() models.Result[any] {
	return i18n.Error(receiver.GetCode(), string(receiver))
}

func (receiver MessageId) MGErrorWithMsg(msg string) models.Result[any] {
	return i18n.ErrorWithMsg(receiver.GetCode(), string(receiver), msg)
}

func (receiver MessageId) MGErrorWithArgs(args ...any) models.Result[any] {
	return i18n.Error(receiver.GetCode(), i18n.Format(string(receiver), args...))
}

func (receiver MessageId) CheckParametersLost(params map[string]string, paramNames ...string) models.Result[any] {
	for _, param := range paramNames {
		v := params[param]
		if v == "" {
			return i18n.Error(receiver.GetCode(), i18n.Format(string(receiver)+":{}", param))
		}
	}
	return models.Success[any](nil)
}

func (receiver MessageId) GetText() string {
	text := i18n.String(string(receiver))
	if text == "" {
		return string(receiver)
	}
	return text
}

func (receiver MessageId) GetCode() int {
	switch receiver {
	case Success:
		return 1
	case UrlNotFound:
		return errcode.URI_NOT_FOUND
	case DbConnectErr:
		return errcode.DB_CONNECT_ERROR
	case ParamLost, ParamError:
		return errcode.REQUEST_PARAMETER_LOST
	case ServiceUnavailable:
		return errcode.SERVICE_UNAVAILABLE
	default:
		return -1
	}
}

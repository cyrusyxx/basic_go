package ioc

import (
	"webook/webook/internal/service/sms"
	"webook/webook/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
}

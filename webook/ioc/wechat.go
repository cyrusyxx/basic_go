package ioc

import (
	"os"
	"webook/webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID is required")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET is required")
	}
	return wechat.NewWechatService(appID, appSecret)
}

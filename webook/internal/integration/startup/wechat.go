package startup

import "webook/webook/internal/service/oauth2/wechat"

func InitWechatService() wechat.Service {
	return wechat.NewWechatService("", "")
}

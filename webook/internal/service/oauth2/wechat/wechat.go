package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"webook/webook/internal/domain"
)

type Service interface {
	AuthUrl(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type WechatService struct {
	appID     string
	appSecret string
	client    *http.Client
}

type Result struct {
	AccessToken  string `json:"access_token"`
	ExpireIn     int64  `json:"expire_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`

	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewWechatService(appID string, appSecret string) Service {
	return &WechatService{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

// https://open.weixin.qq.com/connect/qrconnect?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
func (w *WechatService) AuthUrl(ctx context.Context) (string, error) {
	var redirectUrl = url.PathEscape("https://cyrusss.top/oauth2/wechat/callback")
	const authUrlPattern = `https://open.weixin.qq.com/connect/qrconnect?` +
		`appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	state := "3d6be0a4035d839573b04816624a415e" // TODO state must be random

	return fmt.Sprintf(authUrlPattern, w.appID, redirectUrl, state), nil
}

func (w *WechatService) VerifyCode(ctx context.Context,
	code string) (domain.WechatInfo, error) {
	const curl = `https://api.weixin.qq.com/sns/oauth2/access_token?` +
		`appid=%s&secret=%s&code=%s&grant_type=authorization_code`
	tokenUrl := fmt.Sprintf(curl, w.appID, w.appSecret, code)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := w.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	var res Result
	if err := json.NewDecoder(httpResp.Body).Decode(&res); err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf(
			"errcode: %d, errmsg: %s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		UnionId: res.UnionID,
		OpenId:  res.OpenID,
	}, nil
}

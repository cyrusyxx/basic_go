package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/webook/internal/service"
	"webook/webook/internal/service/oauth2/wechat"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/pkg/ginx"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	usersvc service.UserService
	ijwt.Handler
}

func NewOAuth2WechatHandler(svc wechat.Service,
	usersvc service.UserService, jwthdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		usersvc: usersvc,
		Handler: jwthdl,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2Url)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2Url(ctx *gin.Context) {
	// TODO implement state
	// TODO in 六.4
	panic("implement me")
	url, err := o.svc.AuthUrl(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "generate auth url failed",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: url,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	//state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "verify code failed",
		})
		return
	}
	u, err := o.usersvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "find or create user failed",
		})
		return
	}
	o.SetJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "Login success",
	})
	return
}

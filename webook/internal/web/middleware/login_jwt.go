package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	ijwt "webook/webook/internal/web/jwt"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(hdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip login check for signup and login
		path := ctx.Request.URL.Path
		if path == "/user/signup" ||
			path == "/user/login" ||
			path == "/user/login_sms/code/send" ||
			path == "/user/login_sms/code/verify" ||
			path == "/oauth2/wechat/callback" ||
			path == "/oauth2/wechat/authurl" {
			return
		}

		// Get token
		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims

		// Parse token
		token, err := jwt.ParseWithClaims(tokenStr, &uc,
			func(token *jwt.Token) (interface{}, error) {
				return ijwt.SigKey, nil
			})

		// Check if token is valid
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check ssid
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userclaim", uc)
	}
}

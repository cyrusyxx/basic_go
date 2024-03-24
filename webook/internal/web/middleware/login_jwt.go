package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"webook/webook/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip login check for signup and login
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms/code/verify" ||
			path == "/oauth2/wechat/callback" ||
			path == "/oauth2/wechat/authurl" {
			return
		}

		// Get token
		tokenStr := web.ExtractToken(ctx)
		var uc web.UserClaims

		// Parse token
		token, err := jwt.ParseWithClaims(tokenStr, &uc,
			func(token *jwt.Token) (interface{}, error) {
				return web.SigKey, nil
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

		ctx.Set("userclaim", uc)
	}
}

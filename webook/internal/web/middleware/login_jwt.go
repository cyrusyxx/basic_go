package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/webook/constants"
	"webook/webook/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip login check for signup and login
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}

		// Check if token is valid
		authStr := ctx.GetHeader("Authorization")
		if authStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if token is valid
		segs := strings.Split(authStr, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Get token
		tokenStr := segs[1]
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

		// Check if token is about to expire
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(constants.CheckLoginExpireTime))
			tokenStr, err = token.SignedString(web.SigKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.Set("userclaim", uc)

	}
}

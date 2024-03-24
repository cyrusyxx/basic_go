package jwt

import "github.com/gin-gonic/gin"

type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetJWTToken(ctx *gin.Context, uid int64)
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
}

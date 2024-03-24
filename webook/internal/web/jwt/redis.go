package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
	"webook/webook/constants"
)

var SigKey = []byte("ukRIDSD0JpWD5Qv0P46Y8IGLjB2uvShj")
var RefreshSigKey = []byte("ckRIDSD0JpWD5Qv0P46Y8IGLjB2uvShj")

type RedisJWTHandler struct {
	signingMethod jwt.SigningMethod
	client        redis.Cmdable
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
	Ssid      string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		signingMethod: jwt.SigningMethodHS512,
		client:        cmd,
	}
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64) {
	ssid := uuid.New().String()
	err := h.SetRefreshJWTToken(ctx, uid, ssid)
	if err != nil {
		ctx.String(http.StatusOK, "System Error!!")
		return
	}

	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.JwtExpireTime)),
		},
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		Ssid:      ssid,
	}

	// Generate token
	token := jwt.NewWithClaims(h.signingMethod, uc)
	tokenStr, err := token.SignedString(SigKey)
	if err != nil {
		ctx.String(http.StatusOK, "System Error!!")
		return
	}

	// Set the token to the header
	ctx.Header("x-jwt-token", tokenStr)
}

func (h *RedisJWTHandler) SetRefreshJWTToken(ctx *gin.Context,
	uid int64, ssid string) error {
	uc := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().
				Add(constants.JwtRefreshExpireTime)),
		},
		Uid:  uid,
		Ssid: ssid,
	}

	// Generate token
	token := jwt.NewWithClaims(h.signingMethod, uc)
	tokenStr, err := token.SignedString(RefreshSigKey)
	if err != nil {
		ctx.String(http.StatusOK, "System Error!!")
		return err
	}

	// Set the token to the header
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	claims := ctx.MustGet("userclaim").(UserClaims)
	return h.client.Set(ctx, "users:ssid:"+claims.Ssid, "",
		constants.JwtRefreshExpireTime).Err()
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	// Check ssid
	cnt, err := h.client.Exists(ctx, "users:ssid:"+ssid).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("invalid session")
	}
	return nil
}

func (r *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authStr := ctx.GetHeader("Authorization")
	if authStr == "" {
		return ""
	}
	segs := strings.Split(authStr, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

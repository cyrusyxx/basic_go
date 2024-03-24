package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"webook/webook/constants"
)

var SigKey = []byte("ukRIDSD0JpWD5Qv0P46Y8IGLjB2uvShj")
var RefreshSigKey = []byte("ckRIDSD0JpWD5Qv0P46Y8IGLjB2uvShj")

type jwtHandler struct {
	signingMethod jwt.SigningMethod
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

func newJwtHandler() jwtHandler {
	return jwtHandler{
		signingMethod: jwt.SigningMethodHS512,
	}
}

func ExtractToken(ctx *gin.Context) string {
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

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) {
	err := h.setRefreshJWTToken(ctx, uid)
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

func (h *jwtHandler) setRefreshJWTToken(ctx *gin.Context, uid int64) error {
	uc := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().
				Add(constants.JwtRefreshExpireTime)),
		},
		Uid: uid,
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

package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"webook/webook/pkg/logger"
)

var L logger.Logger = logger.NewNopLogger()

var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}

func WrapBodyAndClaims[Req any, Claims any](
	bizFn func(ctx *gin.Context, req Req, uc Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind request
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("Failed to bind request", logger.Error(err))
			return
		}

		// Get user claims and type assertion
		val, ok := ctx.Get("userclaim")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Call business function
		result, err := bizFn(ctx, req, uc)
		vector.WithLabelValues(strconv.Itoa(result.Code)).Inc()
		if err != nil {
			L.Error("Failed to call business function", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func WrapBody[Req any](
	bizFn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind request
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("Failed to bind request", logger.Error(err))
			return
		}

		// Call business function
		result, err := bizFn(ctx, req)
		vector.WithLabelValues(strconv.Itoa(result.Code)).Inc()
		if err != nil {
			L.Error("Failed to call business function", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func WrapClaims[Claims any](
	bizFn func(ctx *gin.Context, uc Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user claims and type assertion
		val, ok := ctx.Get("userclaim")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Call business function
		result, err := bizFn(ctx, uc)
		vector.WithLabelValues(strconv.Itoa(result.Code)).Inc()
		if err != nil {
			L.Error("Failed to call business function", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, result)
	}
}

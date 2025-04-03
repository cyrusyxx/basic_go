package ioc

import (
	"context"
	"webook/webook/constants"
	"webook/webook/internal/web"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/internal/web/middleware"
	"webook/webook/pkg/ginx"
	myprometheus "webook/webook/pkg/ginx/middleware/prometheus"
	"webook/webook/pkg/ginx/middleware/ratelimit"
	"webook/webook/pkg/limiter"
	"webook/webook/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func InitWebServer(middlewareFuncs []gin.HandlerFunc,
	userHandler *web.UserHandler, wechatHandler *web.OAuth2WechatHandler,
	artiHandler *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewareFuncs...)
	userHandler.RegisterRoutes(server)
	wechatHandler.RegisterRoutes(server)
	artiHandler.RegisterRoutes(server)
	return server
}

func InitMiddleware(redisdb redis.Cmdable, hdl ijwt.Handler, lger logger.Logger) []gin.HandlerFunc {
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "webook",
		Subsystem: "webook_backend",
		Name:      "biz_code",
	})
	return []gin.HandlerFunc{
		// Use Middlewares
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			AllowOrigins: []string{"http://localhost:30001",
				"http://localhost:3000"},
			//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
			AllowHeaders:  []string{"Content-Type", "Authorization"},
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
			//ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			//AllowOriginFunc: func(origin string) bool {
			//	return strings.Contains(origin, "localhost")
			//},
			MaxAge: constants.MaxCorsAge,
		}),
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context,
			al middleware.AccessLog) {
			lger.Debug("AccessLog", logger.Field{
				Key:   "req",
				Value: al,
			})
		}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(
			redisdb, constants.RateLimitInterval, constants.RateLimitRate,
		)).Build(),
		myprometheus.NewBuilder("webook",
			"webook_backend",
			"gin_http",
			"",
		).BuildResponseTime(),
		myprometheus.NewBuilder("webook",
			"webook_backend",
			"gin_http",
			"",
		).BuildActiveRequest(),
		otelgin.Middleware("webook"),
	}
}

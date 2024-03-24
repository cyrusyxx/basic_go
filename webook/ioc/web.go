package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"webook/webook/constants"
	"webook/webook/internal/web"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/internal/web/middleware"
	"webook/webook/pkg/ginx/middleware/ratelimit"
	"webook/webook/pkg/limiter"
)

func InitWebServer(middlewareFuncs []gin.HandlerFunc,
	userHandler *web.UserHandler, wechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewareFuncs...)
	userHandler.RegisterRoutes(server)
	wechatHandler.RegisterRoutes(server)
	return server
}

func InitMiddleware(redisdb redis.Cmdable, hdl ijwt.Handler) []gin.HandlerFunc {
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
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(
			redisdb, constants.RateLimitInterval, constants.RateLimitRate,
		)).Build(),
	}
}

package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"webook/webook/constants"
	"webook/webook/internal/web"
	"webook/webook/internal/web/middleware"
	"webook/webook/pkg/ginx/middleware/ratelimit"
)

func InitWebServer(middlewareFuncs []gin.HandlerFunc,
	userHandler *web.UserHandler) *gin.Engine {

	server := gin.Default()
	server.Use(middlewareFuncs...)
	userHandler.RegisterRoutes(server)
	return server
}

func InitMiddleware(redisdb redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// Use Middlewares
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			AllowOrigins: []string{"http://localhost:30001",
				"http://localhost:3000"},
			//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
			AllowHeaders:  []string{"Content-Type", "Authorization"},
			ExposeHeaders: []string{"x-jwt-token"},
			//ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			//AllowOriginFunc: func(origin string) bool {
			//	return strings.Contains(origin, "localhost")
			//},
			MaxAge: constants.MaxCorsAge,
		}),
		middleware.LoginJWTMiddlewareBuilder{}.CheckLogin(),
		ratelimit.NewBuilder(redisdb, constants.RateLimitInterval, constants.RateLimitRate).Build(),
	}
}

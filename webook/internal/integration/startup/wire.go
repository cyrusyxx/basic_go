//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	user3 "webook/webook/internal/repository"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		InitMysql, InitRedis, ioc.InitSMSService,

		dao.NewGORMUserDAO,

		cache.NewRedisUserCache, cache.NewRedisCodeCache,

		user3.NewCachedUserRepository, user3.NewCachedCodeRepository,

		service.NewCachedCodeService, service.NewCachedUserService,

		user.NewUserHandler,

		ioc.InitWebServer, ioc.InitMiddleware,
	)

	return gin.Default()
}

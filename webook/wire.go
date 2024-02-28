//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
	"webook/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitMysql, ioc.InitRedis, ioc.InitSMSService,

		dao.NewGORMUserDAO,

		cache.NewRedisUserCache, cache.NewRedisCodeCache,

		repository.NewCachedUserRepository, repository.NewCachedCodeRepository,

		service.NewCachedCodeService, service.NewCachedUserService,

		web.NewUserHandler,

		ioc.InitWebServer, ioc.InitMiddleware,
	)

	return gin.Default()
}

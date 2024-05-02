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
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/ioc"
)

var interactiveSet = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitMysql,
		ioc.InitRedis,
		ioc.InitLogger,

		interactiveSet,

		dao.NewGORMUserDAO,
		dao.NewGORMArticleDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository.NewCachedArticleRepository,

		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewCachedCodeService,
		service.NewCachedUserService,
		service.NewImplArticleService,

		ijwt.NewRedisJWTHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,

		ioc.InitWebServer,
		ioc.InitMiddleware,
	)

	return gin.Default()
}

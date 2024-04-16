//go:build wireinject

package startup

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

var thirdPartySet = wire.NewSet(
	InitMysql,
	InitRedis,
	InitLogger,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		ioc.InitSMSService,

		dao.NewGORMUserDAO,
		dao.NewGORMArticleDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		service.NewCachedCodeService,
		service.NewCachedUserService,
		InitWechatService,

		ijwt.NewRedisJWTHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		InitArticleHandler,

		ioc.InitWebServer, ioc.InitMiddleware,
	)

	return gin.Default()
}

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		repository.NewCachedArticleRepository,
		service.NewImplArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}

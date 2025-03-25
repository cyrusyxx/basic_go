//go:build wireinject

package main

import (
	"webook/webook/internal/events/article"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/ioc"

	"time"

	"github.com/google/wire"
)

var interactiveSet = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

var rankingSvcSet = wire.NewSet(
	wire.Value(time.Minute*3), // expiration for ranking cache
	cache.NewRankingLocalCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitMysql,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitSaramaClient,
		ioc.InitSyncProducer,
		ioc.InitConsumers,
		ioc.InitJobs,
		ioc.InitRankingJob,
		ioc.InitRlockClient,

		interactiveSet,
		rankingSvcSet,

		article.NewSaramaSyncProducer,
		article.NewInteractiveReadEventConsumer,

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

		wire.Struct(new(App), "*"),
	)

	return new(App)
}

// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"webook/webook/internal/events/article"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
	"webook/webook/internal/web/jwt"
	"webook/webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := ioc.InitLogger()
	v := ioc.InitMiddleware(cmdable, handler, logger)
	db := ioc.InitMysql(logger)
	userDAO := dao.NewGORMUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewCachedUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCachedCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	articleDAO := dao.NewGORMArticleDAO(db)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, articleCache, userRepository)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewImplArticleService(articleRepository, producer, logger)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	localRankingCache := cache.NewRankingLocalCache()
	redisRankingCache := cache.NewRedisRankingCache(cmdable)
	rankingRepository := repository.NewCachedRankingRepository(localRankingCache, redisRankingCache)
	rankingService := service.NewBatchRankingService(rankingRepository, interactiveService, articleService)
	commentDAO := dao.NewGORMCommentDAO(db)
	commentRepository := repository.NewCommentRepo(commentDAO)
	commentService := service.NewCommentServiceImpl(commentRepository)
	articleHandler := web.NewArticleHandler(logger, articleService, interactiveService, rankingService, commentService)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(interactiveRepository, client, logger)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	rlockClient := ioc.InitRlockClient(cmdable)
	rankingJob := ioc.InitRankingJob(rankingService, rlockClient, logger)
	cron := ioc.InitJobs(logger, rankingJob)
	app := &App{
		server:    engine,
		consumers: v2,
		cron:      cron,
	}
	return app
}

// wire.go:

var interactiveSet = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)

var rankingSvcSet = wire.NewSet(cache.NewRedisRankingCache, cache.NewRankingLocalCache, repository.NewCachedRankingRepository, service.NewBatchRankingService)

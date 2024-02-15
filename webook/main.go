package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"path/filepath"
	"webook/webook/constants"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
	"webook/webook/internal/web/middleware"
	"webook/webook/pkg/ginx/middleware/ratelimit"
)

//none

func main() {
	// Init Server and Database
	//gin.SetMode(gin.DebugMode)
	initConfig()
	server := initServer()
	db := initDB()

	// Init UserHandler
	initUserHandler(db, server)

	// Run Server
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initServer() *gin.Engine {
	server := gin.Default()
	handlecors(server)
	//handleSessions(server)
	handlejwt(server)
	handleRatelimit(server)
	return server
}

func initDB() *gorm.DB {
	// Connect Database
	dsn := viper.GetString("mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Init Tables
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	// RegisterRoutes
	hdl.RegisterRoutes(server)
}

func initConfig() {
	// Init Pflag
	configfile := pflag.StringP("config", "c",
		"config/config.yaml", "config file")
	pflag.Parse()

	// Init Viper
	viper.SetConfigFile(filepath.FromSlash(*configfile))
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file not found")
		}
		if _, ok := err.(viper.ConfigParseError); ok {
			panic("Config file parse error")
		}
		panic(err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		println("Config file changed:",
			e.Name, e.Op)
	})

	viper.WatchConfig()
}

// Middle Were Handler
func handlecors(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
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
	}))
}

func handlejwt(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func handleRatelimit(server *gin.Engine) {
	redisaddr := viper.GetString("redis.addr")
	redisDB := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})
	server.Use(ratelimit.NewBuilder(redisDB, constants.RateLimitInterval, constants.RateLimitRate).Build())
}

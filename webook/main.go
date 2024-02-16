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
	handlejwt(server)
	handleRatelimit(server)

	return server
}

func initDB() *gorm.DB {
	// Get DSN From Viper
	dsn := viper.GetString("mysql.dsn")

	// Connect to Database
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

	// Watch Config
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
	// Get Builder
	login := middleware.LoginJWTMiddlewareBuilder{}

	// Use Middleware
	server.Use(login.CheckLogin())
}

func handleRatelimit(server *gin.Engine) {
	// Get Redis Addr From Viper
	redisaddr := viper.GetString("redis.addr")

	// Init Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})

	// Use Ratelimit Middleware
	server.Use(ratelimit.NewBuilder(redisDB, constants.RateLimitInterval, constants.RateLimitRate).Build())
}

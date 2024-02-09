package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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
	gin.SetMode(gin.DebugMode)
	server := initServer()
	db := initDB()

	// Init UserHandler
	initUserHandler(db, server)

	// Run Server
	server.Run(":8080")
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
	dsn := "root:030208@tcp(cyrusss.top:30997)/test?charset=utf8&parseTime=True&loc=Local"
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

// Middle Were Handler
func handlecors(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: []string{"http://localhost:30001", "http://localhost:3000"},
		//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return strings.Contains(origin, "localhost")
		//},
		MaxAge: 12 * time.Hour,
	}))
}

func handlejwt(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func handleRatelimit(server *gin.Engine) {
	redisDB := redis.NewClient(&redis.Options{
		Addr: "cyrusss.top:32699",
	})
	server.Use(ratelimit.NewBuilder(redisDB, time.Minute, 100).Build())
}

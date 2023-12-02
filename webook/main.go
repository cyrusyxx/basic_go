package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
)
//none

func main() {
	// Init Server and Database
	server := initServer()
	db := initDB()

	// Init Userhandler
	initUserHandler(db, server)

	// Run Server
	server.Run(":8080")
}

func initServer() *gin.Engine {
	server := gin.Default()
	handlecors(server)
	return server
}

func initDB() *gorm.DB {
	// Connect Database
	dsn := "root:030208@tcp(cyruss.cn:3306)/test?charset=utf8&parseTime=True&loc=Local"
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
		AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Content-Type"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return strings.Contains(origin, "localhost")
		//},
		MaxAge: 12 * time.Hour,
	}))
}

/*
// JSON 数据
var json_data = {
email: "9328123@qq.com",
password: "Cc@002300",
confirmPassword: "Cc@002300"
};

// 发送 POST 请求
fetch('http://localhost:8080/users/signup', {
method: 'POST',
headers: {
'Content-Type': 'application/json'
},
body: JSON.stringify(json_data)
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error(error));
*/

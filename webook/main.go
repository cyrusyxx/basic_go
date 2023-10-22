package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"webook/webook/internal/web"
)

func main() {
	server := gin.Default()
	handlecors(server)

	hdl := web.NewUserHandler()
	hdl.RegisterRoutes(server)

	server.Run(":8080")
}

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

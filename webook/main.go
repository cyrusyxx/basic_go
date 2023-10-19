package main

import (
	"github.com/gin-gonic/gin"
	"webook/webook/internal/web"
)

func main() {
	server := gin.Default()

	hdl := &web.UserHandler{}
	hdl.RegisterRoutes(server)

	server.Run(":8080")
}

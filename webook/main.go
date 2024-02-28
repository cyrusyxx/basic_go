package main

import (
	"webook/webook/config"
)

//none

func main() {
	config.InitConfig()
	server := InitWebServer()

	// Run Server
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

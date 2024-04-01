package main

import (
	"go.uber.org/zap"
	"webook/webook/config"
)

//none

func main() {
	config.InitConfig()
	initLogger()
	//Testlog()
	server := InitWebServer()

	// Run Server
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func Testlog() {
	type st struct {
		Name string
		Age  int
	}
	t := st{
		Name: "test",
		Age:  18,
	}
	zap.L().Info("Test log", zap.Any("teststruct", t))
	zap.L().Warn("Test log warn")
	zap.L().Error("Test log error")
	zap.L().Debug("Test log debug")
}

package ioc

import (
	"go.uber.org/zap"
	"webook/webook/pkg/logger"
)

func InitLogger() logger.Logger {
	//cfg := zap.NewDevelopmentConfig()
	//err := viper.Unmarshal("log", &cfg)
	//if err != nil {
	//	panic(err)
	//}
	//zaplger, err := cfg.Build()
	zaplger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(zaplger)
}

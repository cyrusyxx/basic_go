package startup

import "webook/webook/pkg/logger"

func InitLogger() logger.Logger {
	return logger.NewNopLogger()
}

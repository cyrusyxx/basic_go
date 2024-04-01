package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"webook/webook/internal/repository/dao"
	"webook/webook/pkg/logger"
)

func InitMysql(l logger.Logger) *gorm.DB {
	// Get DSN From Viper
	dsn := viper.GetString("mysql.dsn")

	// Connect to Database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glogger.New(gormloggerFunc(l.Debug), glogger.Config{
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})
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

type gormloggerFunc func(msg string, fields ...logger.Field)

func (g gormloggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{
		Key:   "args",
		Value: i,
	})
}

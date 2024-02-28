package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/webook/internal/repository/dao"
)

func InitMysql() *gorm.DB {
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

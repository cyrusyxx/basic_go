package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/webook/internal/repository/dao"
)

func InitMysql() *gorm.DB {
	// Connect to Database
	dsn := `root:030208@tcp(cyrusss.top:30997)/test?charset=utf8&parseTime=True&loc=Local`
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

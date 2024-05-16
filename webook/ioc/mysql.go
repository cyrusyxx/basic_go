package ioc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	prometheusp "gorm.io/plugin/prometheus"
	"webook/webook/internal/repository/dao"
	"webook/webook/pkg/gormx"
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

	// Use Prometheus
	err = db.Use(prometheusp.New(prometheusp.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheusp.MetricsCollector{
			&prometheusp.MySQL{},
		},
	}))
	if err != nil {
		panic(err)
	}

	// Use Callback
	cbks := gormx.NewCallbacks(prometheus.SummaryOpts{
		Namespace: "webook",
		Subsystem: "gorm",
		Name:      "gorm_db",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	err = db.Callback().Create().Before("*").
		Register("prometheus_gorm_create_before", cbks.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_gorm_create_after", cbks.After("create"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("*").
		Register("prometheus_gorm_update_before", cbks.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_gorm_update_after", cbks.After("update"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Before("*").
		Register("prometheus_gorm_delete_before", cbks.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").
		Register("prometheus_gorm_delete_after", cbks.After("delete"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Query().Before("*").
		Register("prometheus_gorm_query_before", cbks.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Query().After("*").
		Register("prometheus_gorm_query_after", cbks.After("query"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().Before("*").
		Register("prometheus_gorm_row_before", cbks.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_gorm_row_after", cbks.After("row"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().Before("*").
		Register("prometheus_gorm_raw_before", cbks.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").
		Register("prometheus_gorm_raw_after", cbks.After("raw"))
	if err != nil {
		panic(err)
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

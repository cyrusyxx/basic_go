package main

import (
	"net/http"
	"webook/webook/config"
	"webook/webook/ioc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

//none

func main() {
	// Init CONFIG, LOGGER, APP, PROMETHEUS
	config.InitConfig()
	initLogger()
	app := InitWebServer()
	initPrometheus()

	// kafka health
	if err := ioc.KafkaHealthCheck([]string{"cyruss.cn:9092"}); err != nil {
		panic(err)
	}

	// Start consumers
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	// Start Cron Job
	app.cron.Start()
	defer app.cron.Stop()

	// Run Server
	server := app.server
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

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()
}

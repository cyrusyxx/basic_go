package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"webook/webook/config"
)

//none

func main() {
	// Init CONFIG, LOGGER, APP, PROMETHEUS
	config.InitConfig()
	initLogger()
	app := InitWebServer()
	initPrometheus()

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

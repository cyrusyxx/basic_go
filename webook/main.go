package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"webook/webook/config"
)

//none

func main() {
	config.InitConfig()
	initLogger()
	//Testlog()
	app := InitWebServer()
	initPrometheus()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.server

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

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()
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

//
//func mainSarama() {
//	//	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
//	var kafkaAddr = []string{"cyrusss.top:9094"}
//	cfg := sarama.NewConfig()
//	cfg.Producer.Return.Successes = true
//	producer, err := sarama.NewSyncProducer(kafkaAddr, cfg)
//	if err != nil {
//		panic(err)
//	}
//
//	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
//		Topic: "test_topic",
//		Value: sarama.StringEncoder("test:test"),
//		Headers: []sarama.RecordHeader{
//			{
//				Key:   []byte("key11"),
//				Value: []byte("value11"),
//			},
//		},
//		Metadata: "this is metadata",
//	})
//	if err != nil {
//		panic(err)
//	}
//}

package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
	"webook/webook/config"
	"webook/webook/internal/repository/dao"
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

func mainmongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Set Monitor
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println("command started")
		},
		Succeeded: nil,
		Failed:    nil,
	}

	// Connect mongoDB
	opts := options.Client().
		ApplyURI("mongodb://root:030208@cyrusss.top:31627/").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		print(err)
	}

	// Insert One
	col := client.Database("webook").Collection("articles")
	insertRes, err := col.InsertOne(ctx, dao.Article{
		Id:       3,
		Title:    "My Title",
		Content:  "My content",
		AuthorId: 123,
		States:   0,
		Ctime:    0,
		Utime:    0,
	})
	if err != nil {
		fmt.Println(err)
	}

	// Print OID
	oid := insertRes.InsertedID.(primitive.ObjectID)
	fmt.Println(oid)

	// Find col
	filter := bson.M{
		"id": 3,
	}
	findRes := col.FindOne(ctx, filter)

	// Decode to arti
	var arti dao.Article
	err = findRes.Decode(&arti)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", arti)

	fmt.Println(err)
}

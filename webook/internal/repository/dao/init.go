package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Article{},
		&PublicArticle{},
		&InteractiveCount{},
		&UserLikeBiz{},
		&UserCollectionBiz{},
		&Comment{},
	)
}

func InitCollection(mdb *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := mdb.Collection("articles")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"author_id", 1}},
		},
	})
	if err != nil {
		return err
	}
	livecol := mdb.Collection("published_articles")
	_, err = livecol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"author_id", 1}},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

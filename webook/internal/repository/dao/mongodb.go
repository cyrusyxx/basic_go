package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBArticleDAO struct {
	node    *snowflake.Node
	col     *mongo.Collection
	liveCol *mongo.Collection
}

func NewMongoDBArticleDAO(db *mongo.Database,
	node *snowflake.Node) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		node:    node,
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
	}

}

func (d *MongoDBArticleDAO) Insert(ctx context.Context,
	arti Article) (int64, error) {
	now := time.Now().UnixMilli()
	arti.Ctime = now
	arti.Utime = now
	arti.Id = d.node.Generate().Int64()

	// Insert
	_, err := d.col.InsertOne(ctx, &arti)
	if err != nil {
		return 0, err
	}
	return arti.Id, nil
}

func (d *MongoDBArticleDAO) UpdateById(ctx context.Context,
	arti Article) error {
	filter := bson.M{
		"id":        arti.Id,
		"author_id": arti.AuthorId,
	}
	set := bson.D{
		{"$set", bson.M{
			"title":   arti.Title,
			"content": arti.Content,
			"utime":   time.Now().UnixMilli(),
			"status":  arti.Status},
		},
	}
	res, err := d.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("id or author_id is wrong")
	}
	return nil
}

func (d *MongoDBArticleDAO) Sync(ctx context.Context,
	arti Article) (int64, error) {
	var (
		id  = arti.Id
		err error
	)
	if id > 0 {
		err = d.UpdateById(ctx, arti)
	} else {
		id, err = d.Insert(ctx, arti)
	}
	if err != nil {
		return 0, err
	}

	arti.Id = id
	arti.Utime = time.Now().UnixMilli()
	filter := bson.D{
		{"id", arti.Id},
		{"author_id", arti.AuthorId},
	}
	set := bson.D{
		{"$set", arti},
		{"$setOnInsert", bson.M{"ctime": time.Now().UnixMilli()}},
	}
	_, err = d.liveCol.UpdateOne(ctx, filter, set,
		options.Update().SetUpsert(true))
	return id, err
}

func (d *MongoDBArticleDAO) SyncStatus(ctx context.Context,
	uid int64, id int64, status uint8) error {
	filter := bson.D{
		{"id", id},
		{"author_id", uid},
	}
	set := bson.D{
		{"$set", bson.D{{"status", status}}},
	}

	// Update
	res, err := d.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return errors.New("id or author_id is wrong")
	}
	_, err = d.liveCol.UpdateOne(ctx, filter, set)

	return err
}

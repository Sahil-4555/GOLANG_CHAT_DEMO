package database

import (
	"context"
	"sync"

	"github.com/Sahil-4555/mvc/configs"
	"github.com/Sahil-4555/mvc/models"
	"github.com/Sahil-4555/mvc/shared/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IConnection interface {
	GetMongoDB() *mongo.Client
	UserCollection() *mongo.Collection
	MessageCollection() *mongo.Collection
	ChannelCollection() *mongo.Collection
}

type Connection struct {
	Mongodb *mongo.Client
}

var database = ""
var client *mongo.Client
var ctx = context.TODO()
var connectionOnce sync.Once

func createIndexes(client *mongo.Client) error {
	database := client.Database(configs.Database())

	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"user_name": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"favourites": 1},
		},
		{
			Keys: bson.M{"deleted_at": 1},
		},
		{
			Keys: bson.M{"created_at": 1},
		},
		{
			Keys: bson.M{"updated_at": 1},
		},
	}

	_, err := database.Collection(models.UserCollection).Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return err
	}

	channelIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"channel_type": 1},
		},
		{
			Keys: bson.M{"users": 1},
		},
		{
			Keys: bson.M{"last_activity": 1},
		},
		{
			Keys: bson.M{"last_opened": 1},
		},
		{
			Keys: bson.M{"deleted_at": 1},
		},
		{
			Keys: bson.M{"created_at": 1},
		},
		{
			Keys: bson.M{"updated_at": 1},
		},
	}

	_, err = database.Collection(models.ChannelCollection).Indexes().CreateMany(ctx, channelIndexes)
	if err != nil {
		return err
	}

	messageIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"channel_id": 1},
		},
		{
			Keys: bson.M{"sender": 1},
		},
		{
			Keys: bson.M{"read_by": 1},
		},
		{
			Keys: bson.M{"deleted_at": 1},
		},
		{
			Keys: bson.M{"created_at": 1},
		},
		{
			Keys: bson.M{"updated_at": 1},
		},
	}

	_, err = database.Collection(models.MessageCollection).Indexes().CreateMany(ctx, messageIndexes)
	if err != nil {
		return err
	}

	return nil
}

func Init() {
	database = configs.Database()
	connectionOnce.Do(func() {
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(configs.MongoURI()))
		if err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			panic(err.Error())
		}
		client = mongoClient
		err = client.Ping(ctx, nil)
		if err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			panic(err.Error())
		}

		err = createIndexes(client)
		if err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			panic(err.Error())
		}
	})
}

func NewConnection() *Connection {
	return &Connection{
		client,
	}
}

func Close() error {
	if client != nil {
		if err := client.Disconnect(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (conn *Connection) GetMongoDB() *mongo.Client {
	return conn.Mongodb
}

func (conn *Connection) UserCollection() *mongo.Collection {
	return conn.Mongodb.Database(database).Collection(models.UserCollection)
}

func (conn *Connection) MessageCollection() *mongo.Collection {
	return conn.Mongodb.Database(database).Collection(models.MessageCollection)
}

func (conn *Connection) ChannelCollection() *mongo.Collection {
	return conn.Mongodb.Database(database).Collection(models.ChannelCollection)
}

package data

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	Client *mongo.Client
}

func CreateDBClient(uri string) (*DB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	return &DB{Client: client}, err
}

func PingDB(client *mongo.Client) error {
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return err
	}
	return nil
}

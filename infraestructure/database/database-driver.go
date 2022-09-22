package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateDBClient(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	return client, err
}

func PingDB(client *mongo.Client) error {
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return err
	}
	return nil
}
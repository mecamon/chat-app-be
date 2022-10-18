package repositories_impl

import (
	"context"
	"errors"
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ClusterMsgRepo struct {
	app    *config.App
	dbConn *data.DB
}

var clusterMsgRepo repositories.ClusterMsgRepo

func InitClusterMsgRepo(app *config.App, dbConn *data.DB) repositories.ClusterMsgRepo {
	clusterMsgRepo = &ClusterMsgRepo{
		app:    app,
		dbConn: dbConn,
	}
	return clusterMsgRepo
}

func GetClusterMsgRepo() repositories.ClusterMsgRepo {
	return clusterMsgRepo
}

func (r *ClusterMsgRepo) Create(cluster models.ClusterOfMessages) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cColl := r.dbConn.Client.Database(r.app.DBName).Collection("cluster_of_messages")
	result, err := cColl.InsertOne(ctx, cluster)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *ClusterMsgRepo) Update(clusterID string, message models.MsgContent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	clusterObjectID, err := primitive.ObjectIDFromHex(clusterID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", clusterObjectID}}
	update := bson.D{{"$push", bson.D{{"messages", message}}}}

	cColl := r.dbConn.Client.Database(r.app.DBName).Collection("cluster_of_messages")
	result, err := cColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no document was updated")
	}

	return nil
}

func (r *ClusterMsgRepo) GetLatest() (models.ClusterOfMessages, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	filter := bson.D{}
	opt := &options.FindOneOptions{
		Sort: bson.D{{"_id", -1}},
	}

	var cluster models.ClusterOfMessages

	cColl := r.dbConn.Client.Database(r.app.DBName).Collection("cluster_of_messages")
	result := cColl.FindOne(ctx, filter, opt)
	if result.Err() != nil {
		return cluster, result.Err()
	}

	if err := result.Decode(&cluster); err != nil {
		return cluster, err
	}

	return cluster, nil
}

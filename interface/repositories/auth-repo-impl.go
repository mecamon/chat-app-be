package repositories_impl

import (
	"context"
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AuthRepoImpl struct {
	App    *config.App
	DBConn *data.DB
}

var authRepoImpl *AuthRepoImpl

func InitAuthRepo(app *config.App, dbConn *data.DB) repositories.AuthRepo {
	authRepoImpl = &AuthRepoImpl{
		App:    app,
		DBConn: dbConn,
	}
	return authRepoImpl
}

func GetAuthRepo() repositories.AuthRepo {
	return authRepoImpl
}

func (a *AuthRepoImpl) Register(uEntry models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uColl := a.DBConn.Client.Database(a.App.DBName).Collection("users")
	result, err := uColl.InsertOne(ctx, uEntry)
	if err != nil {
		return "", err
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	return insertedID, nil
}

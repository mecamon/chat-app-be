package repositories_impl

import (
	"context"
	"errors"
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"
	"go.mongodb.org/mongo-driver/bson"
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

func (a *AuthRepoImpl) Login(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	uColl := a.DBConn.Client.Database(a.App.DBName).Collection("users")
	filter := bson.D{{"email", email}}
	result := uColl.FindOne(ctx, filter)

	if result.Err() != nil {
		return "", errors.New("wrong email or password")
	}
	if err := result.Decode(&user); err != nil {
		return "", err
	}

	hasCorrectPass, err := utils.CompareHashAndPass(user.Password, password)
	if err != nil {
		return "", err
	}
	if !hasCorrectPass {
		return "", errors.New("wrong email or password")
	}
	ID := user.ID.Hex()
	return ID, nil
}

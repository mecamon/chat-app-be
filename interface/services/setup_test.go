//go:build integration
// +build integration

package services

import (
	"context"
	"fmt"
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	repositories_impl "github.com/mecamon/chat-app-be/interface/repositories"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
)

var app *config.App
var authRepo repositories.AuthRepo
var authTestService *Auth
var mailTestService *Mail

func TestMain(m *testing.M) {
	dbConn := runDB()
	runRepos(dbConn)
	runServices(app, repositories_impl.GetAuthRepo())
	code := m.Run()
	shutdown(dbConn)
	os.Exit(code)
}

func runDB() *data.DB {
	config.SetConfig()
	app = config.GetConfig()

	dbConnUri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?maxPoolSize=20&w=majority",
		app.DBUser,
		app.DBUserPassword,
		app.DBHost,
		app.DBPort,
		app.DBName)
	dbConn, err := data.CreateDBClient(dbConnUri)
	if err != nil {
		panic(err.Error())
	}
	return dbConn
}

func runRepos(dbConn *data.DB) {
	authRepo = repositories_impl.InitAuthRepo(app, dbConn)
}

func runServices(app *config.App, repo repositories.AuthRepo) {
	authTestService = InitAuth(app, repo)
	mailTestService = InitMailService(app)
}

func shutdown(dbConn *data.DB) {
	uColl := dbConn.Client.Database(app.DBName).Collection("users")
	_, err := uColl.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err.Error())
	}

	if err := dbConn.Client.Disconnect(context.TODO()); err != nil {
		log.Println(err.Error())
	}
}

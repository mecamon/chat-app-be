package repositories_impl

import (
	"context"
	"fmt"
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
)

var dbConn *data.DB
var authTestRepo repositories.AuthRepo
var groupChatTestRepo repositories.GroupChat
var clusterTestMsgRepo repositories.ClusterMsgRepo
var app *config.App

func TestMain(m *testing.M) {
	dbConn = run()
	code := m.Run()
	shutdown(dbConn)
	os.Exit(code)
}

func run() *data.DB {
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
	authTestRepo = InitAuthRepo(app, dbConn)
	groupChatTestRepo = InitGroupChatRepo(app, dbConn)
	clusterTestMsgRepo = InitClusterMsgRepo(app, dbConn)
	return dbConn
}

func shutdown(dbConn *data.DB) {
	uColl := dbConn.Client.Database(app.DBName).Collection("users")
	_, err := uColl.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err.Error())
	}

	gColl := dbConn.Client.Database(app.DBName).Collection("chat_groups")
	_, err = gColl.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err.Error())
	}

	if err := dbConn.Client.Disconnect(context.TODO()); err != nil {
		log.Println(err.Error())
	}
}

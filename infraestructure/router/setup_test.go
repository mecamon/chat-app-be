package router

import (
	"context"
	"fmt"
	"github.com/mecamon/chat-app-be/config"
	appi18n "github.com/mecamon/chat-app-be/i18n"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/interface/controller"
	repositories_impl "github.com/mecamon/chat-app-be/interface/repositories"
	"github.com/mecamon/chat-app-be/interface/services"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"testing"
)

var app *config.App
var mainRouter http.Handler

func TestMain(m *testing.M) {
	config.SetConfig()
	app = config.GetConfig()
	if err := appi18n.InitLocales(); err != nil {
		panic(err.Error())
	}
	mLoc := appi18n.GetMultiLocales()

	dbConn := runDB()
	authTestRepo := repositories_impl.InitAuthRepo(app, dbConn)
	authTestService := services.InitAuth(app, authTestRepo)
	controller.InitAuthController(app, mLoc, authTestService)

	runRouter()
	code := m.Run()
	shutdown(dbConn)
	os.Exit(code)
}

func runDB() *data.DB {
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

func runRouter() {
	SetRouter()
	main, err := GetMain()
	if err != nil {
		panic(err.Error())
	}
	main.AddSubRouters()
	mainRouter = main.R
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

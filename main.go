package main

import (
	"context"
	"fmt"
	appi18n "github.com/mecamon/chat-app-be/i18n"
	"github.com/mecamon/chat-app-be/interface/controller"
	repository "github.com/mecamon/chat-app-be/interface/repositories"
	"github.com/mecamon/chat-app-be/interface/services"
	"log"
	"net/http"
	"time"

	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/infraestructure/router"
)

func main() {

	config.SetConfig()
	app := config.GetConfig()

	err := appi18n.InitLocales()
	if err != nil {
		panic(err)
	}

	dbConn := runDB(app)
	defer func() {
		if err := dbConn.Client.Connect(context.TODO()); err != nil {
			panic(err.Error())
		}
	}()
	runRepos(app, dbConn)

	hub := &services.Hub{
		Clients:        make(map[*services.Client]bool),
		Broadcast:      make(chan services.MessageStruct),
		Register:       make(chan *services.Client),
		Unregister:     make(chan *services.Client),
		AuthRepo:       repository.GetAuthRepo(),
		GroupChatRepo:  repository.GetGroupChatRepo(),
		ClusterMsgRepo: repository.GetClusterMsgRepo(),
	}

	go hub.Run()

	runControllers(hub)
	services.InitMailService(app)
	handler := runRouters()

	srv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("0.0.0.0%s", app.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on port%s...", app.Port)
	log.Fatal(srv.ListenAndServe())
}

func runDB(app *config.App) *data.DB {
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
	if err := data.PingDB(dbConn.Client); err != nil {
		panic(err.Error())
	}
	log.Println("Connected to database...")

	return dbConn
}

func runRepos(app *config.App, dbConn *data.DB) {
	_ = repository.InitAuthRepo(app, dbConn)
	_ = repository.InitGroupChatRepo(app, dbConn)
	_ = repository.InitClusterMsgRepo(app, dbConn)
}

func runControllers(hub *services.Hub) {
	_ = controller.InitAuthController()
	_ = controller.InitGroupChats()
	_ = controller.InitWSHandshake(hub)
}

func runRouters() http.Handler {
	router.SetRouter()
	main, err := router.GetMain()
	if err != nil {
		panic(err.Error())
	}
	main.AddSubRouters()
	return main.R
}

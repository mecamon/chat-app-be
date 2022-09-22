package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/database"
	"github.com/mecamon/chat-app-be/infraestructure/router"
)

func main() {
	config.SetConfig()
	app := config.GetConfig()

	dbConnUri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/test?maxPoolSize=20&w=majority", 
		app.DBUser, 
		app.DBUserPassword, 
		app.DBHost, 
		app.DBPort)
		
	client, err := database.CreateDBClient(dbConnUri)
	if err != nil {
		panic(err.Error())
	}
	if err := database.PingDB(client); err != nil {
		panic(err.Error())
	}
	log.Println("Connected to database...")

	router.SetRouter()
	main, err := router.GetMain()
	if err != nil {
		panic(err.Error())
	}
	main.AddSubRouters()

	srv := &http.Server{
        Handler:      main.R,
		Addr: fmt.Sprintf("0.0.0.0%s", app.Port),
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

	log.Printf("Listening on port%s...", app.Port)
	log.Fatal(srv.ListenAndServe())
}
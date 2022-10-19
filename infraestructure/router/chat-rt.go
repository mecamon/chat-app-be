package router

import (
	"github.com/mecamon/chat-app-be/interface/controller"
	"github.com/mecamon/chat-app-be/interface/middlewares"
	"log"
)

func AddHubSubRouter() {
	main, err := GetMain()
	if err != nil {
		log.Println(err.Error())
	}

	hubController := controller.GetWSHandShake()
	s := main.R.PathPrefix("/api/chat").Subrouter()
	s.HandleFunc("/handshake", hubController.Connect)
	s.Use(middlewares.TokenValidation)
}

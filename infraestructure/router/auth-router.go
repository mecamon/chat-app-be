package router

import (
	"log"

	"github.com/mecamon/chat-app-be/interface/controller"
)

func AddAuthSubrouter() {
	main, err := GetMain()
	if err != nil {
		log.Println(err.Error())
	}
	
	authController := controller.GetAuthController()
	s := main.R.PathPrefix("/api/auth").Subrouter()
	s.HandleFunc("/register", authController.Register)
	s.HandleFunc("/login", authController.Login)
}
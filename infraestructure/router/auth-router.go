package router

import (
	"log"
	"net/http"

	"github.com/mecamon/chat-app-be/interface/controller"
)

func AddAuthSubRouter() {
	main, err := GetMain()
	if err != nil {
		log.Println(err.Error())
	}

	authController := controller.GetAuthController()
	s := main.R.PathPrefix("/api/auth").Subrouter()
	s.HandleFunc("/register", authController.Register).Methods(http.MethodPost)
	s.HandleFunc("/login", authController.Login).Methods(http.MethodPost)
	s.HandleFunc("/recover", authController.SendRecoveryLink).Methods(http.MethodPost)
	s.HandleFunc("/change-password", authController.ChangePass).Methods(http.MethodPost)
}

package router

import (
	"github.com/mecamon/chat-app-be/interface/middlewares"
	"log"
	"net/http"

	"github.com/mecamon/chat-app-be/interface/controller"
)

func AddChatGroupsSubRouter() {
	main, err := GetMain()
	if err != nil {
		log.Println(err.Error())
	}

	chatGroupsController := controller.GetGroupChats()
	s := main.R.PathPrefix("/api/chat_groups").Subrouter()
	s.HandleFunc("/", chatGroupsController.Create).Methods(http.MethodPost)
	s.HandleFunc("/", chatGroupsController.Update).Methods(http.MethodPatch)
	s.HandleFunc("/{group_id}", chatGroupsController.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/add_to_chat/{group_id}", chatGroupsController.AddUserToChat).Methods(http.MethodPost)
	s.HandleFunc("/load_chats", chatGroupsController.LoadAll).Methods(http.MethodGet)
	s.Use(middlewares.TokenValidation)
}

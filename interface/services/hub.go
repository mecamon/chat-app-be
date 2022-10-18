package services

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mecamon/chat-app-be/models"
	"log"
)

type MessageStruct struct {
	MessageType int
	P           []byte
}

type Client struct {
	Hub           *Hub
	Conn          *websocket.Conn
	UserInfo      models.User
	UserGroupsIDs []string
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan MessageStruct
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("New client registered!!!")
		case client := <-h.Unregister:
			delete(h.Clients, client)
			log.Println("Client removed!!!")
		case m := <-h.Broadcast:
			for client := range h.Clients {
				var msgContentDto models.MsgContentDTO
				err := json.Unmarshal(m.P, &msgContentDto)
				if err != nil {
					log.Println(err.Error())
				}

				err = client.Conn.WriteMessage(m.MessageType, m.P)
				if err != nil {
					log.Println(err.Error())
				}

				//if h.isInTheGroup(msgContentDto, client) {
				//	if err != nil {
				//		log.Println(err.Error())
				//	}
				//	err = client.Conn.WriteMessage(m.MessageType, m.P)
				//	if err != nil {
				//		log.Println(err.Error())
				//	}
				//}
			}
		}
	}
}

func (h *Hub) isInTheGroup(m models.MsgContentDTO, client *Client) bool {
	for _, gID := range client.UserGroupsIDs {
		if m.To == gID {
			return true
		}
	}
	return false
}

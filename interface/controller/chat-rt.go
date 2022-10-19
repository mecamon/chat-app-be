package controller

import (
	"github.com/gorilla/websocket"
	repositories_impl "github.com/mecamon/chat-app-be/interface/repositories"
	"github.com/mecamon/chat-app-be/interface/services"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"

	"log"
	"net/http"
)

var WS *WSHandshake

type WSHandshake struct {
	hub           *services.Hub
	authRepo      repositories.AuthRepo
	groupChatRepo repositories.GroupChat
}

func InitWSHandshake(h *services.Hub) *WSHandshake {
	WS = &WSHandshake{
		hub:           h,
		authRepo:      repositories_impl.GetAuthRepo(),
		groupChatRepo: repositories_impl.GetGroupChatRepo(),
	}
	return WS
}

func GetWSHandShake() *WSHandshake {
	return WS
}

func (c *WSHandshake) Connect(w http.ResponseWriter, r *http.Request) {
	ID := r.Context().Value("ID").(string)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := c.authRepo.FindByID(ID)
	if err != nil {
		conn.Close()
		_ = utils.JSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	var skip int64 = 0
	var take int64 = 100

	groups, err := c.groupChatRepo.LoadAll(
		ID,
		map[string]interface{}{
			"chats": "participating",
			"skip":  skip,
			"take":  take,
		})
	if err != nil {
		conn.Close()
		_ = utils.JSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	var groupsID []string

	for _, g := range groups {
		groupsID = append(groupsID, g.ID.Hex())
	}

	client := &services.Client{
		Hub:           c.hub,
		Conn:          conn,
		UserInfo:      user,
		UserGroupsIDs: groupsID,
	}

	c.hub.Register <- client
	defer func() {
		c.hub.Unregister <- client
		conn.Close()
	}()

	for {
		//Read from client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		client.Hub.Broadcast <- services.MessageStruct{
			MessageType: messageType,
			P:           p,
		}
	}
}

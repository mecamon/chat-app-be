//go:build integration
// +build integration

package services

import (
	"github.com/gorilla/websocket"
	repositories_impl "github.com/mecamon/chat-app-be/interface/repositories"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHub_Run(t *testing.T) {
	hub := &Hub{
		Clients:        make(map[*Client]bool),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan MessageStruct),
		AuthRepo:       repositories_impl.GetAuthRepo(),
		GroupChatRepo:  repositories_impl.GetGroupChatRepo(),
		ClusterMsgRepo: repositories_impl.GetClusterMsgRepo(),
	}

	go hub.Run()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "ws://localhost:8080/ws", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-Websocket-Extensions", "permessage-deflate")
	req.Header.Set("Sec-Websocket-Extensions", "client_max_window_bits")
	req.Header.Set("Sec-Websocket-Key", "nt5Yp0gaMC5ieKYpdtP6cA==")
	req.Header.Set("Sec-Websocket-Version", "13")
	req.Header.Set("Upgrade", "websocket")

	conn, err := upgrader.Upgrade(rr, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	u := models.User{
		Name:      "Hub com",
		Bio:       "This is a hub user",
		Email:     "hub@service.com",
		Password:  "validPass123",
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, err := utils.GenerateHash(u.Password)
	if err != nil {
		t.Error(err.Error())
	}
	u.Password = hashedPass

	insertedUID, err := authRepo.Register(u)
	if err != nil {
		t.Error(err.Error())
	}

	userInfo, err := authRepo.FindByID(insertedUID)
	if err != nil {
		t.Error(err.Error())
	}

	objectUID, err := primitive.ObjectIDFromHex(insertedUID)
	group := models.GroupChat{
		Name:        "test group hub",
		Description: "test group hub",
		GroupOwner:  objectUID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	groupID, err := groupChatRepo.Create(insertedUID, group)
	if err != nil {
		t.Error(err.Error())
	}

	client := &Client{
		Hub:           hub,
		Conn:          conn,
		UserInfo:      userInfo,
		UserGroupsIDs: []string{groupID},
	}
	hub.Register <- client

	if len(hub.Clients) != 1 {
		t.Error("expected client length is 1 but did NOT get that")
	}

	hub.Unregister <- client
	if len(hub.Clients) != 0 {
		t.Error("expected client length is 0 but did NOT get that")
	}

	defer func() {
		conn.Close()
		hub.Unregister <- client
	}()
}

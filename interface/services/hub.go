package services

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const limitByCluster = 100

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
	Clients        map[*Client]bool
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan MessageStruct
	GroupAddition  chan models.MsgContentDTO
	AuthRepo       repositories.AuthRepo
	GroupChatRepo  repositories.GroupChat
	ClusterMsgRepo repositories.ClusterMsgRepo
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			delete(h.Clients, client)
		case m := <-h.Broadcast:
			h.handleBroadcast(m)
		case m := <-h.GroupAddition:
			h.groupAddition(m)
		}
	}
}

func (h *Hub) groupAddition(m models.MsgContentDTO) {
	u, err := h.AuthRepo.FindByID(m.From)
	if err != nil {
		log.Println(err.Error())
	}

	err = h.GroupChatRepo.AddUserToChat(u, m.To)
	if err != nil {
		log.Println(err.Error())
	}

	for client := range h.Clients {
		if client.UserInfo.ID.Hex() == u.ID.Hex() {
			client.UserGroupsIDs = append(client.UserGroupsIDs, m.To)
		}
	}
}

func (h *Hub) handleBroadcast(m MessageStruct) {
	var msgContentDto models.MsgContentDTO
	err := json.Unmarshal(m.P, &msgContentDto)
	if err != nil {
		log.Println(err.Error())
	}
	err = h.handleMsgStorage(msgContentDto)
	if err != nil {
		log.Println(err.Error())
	}

	for client := range h.Clients {
		if h.isInTheGroup(msgContentDto, client) {
			err = client.Conn.WriteMessage(m.MessageType, m.P)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (h *Hub) handleMsgStorage(m models.MsgContentDTO) error {
	activeCluster, err := h.ClusterMsgRepo.GetLatest()
	clusterID := activeCluster.ID.Hex()

	if err != nil || len(activeCluster.Messages) == limitByCluster {
		belongsTo, err := primitive.ObjectIDFromHex(m.From)
		if err != nil {
			return err
		}
		ID, err := h.ClusterMsgRepo.Create(models.ClusterOfMessages{
			BelongsToGroup: belongsTo,
			Messages:       nil,
			CreatedAt:      time.Now().Unix(),
			UpdatedAt:      time.Now().Unix(),
			IsClosed:       false,
		})
		clusterID = ID
	}

	from, err := primitive.ObjectIDFromHex(m.From)
	if err != nil {
		return err
	}
	to, err := primitive.ObjectIDFromHex(m.To)
	if err != nil {
		return err
	}

	message := models.MsgContent{
		From:        from,
		To:          to,
		TextContent: m.TextContent,
	}
	err = h.ClusterMsgRepo.Update(clusterID, message)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hub) isInTheGroup(m models.MsgContentDTO, client *Client) bool {
	for _, gID := range client.UserGroupsIDs {
		if m.To == gID {
			return true
		}
	}
	return false
}

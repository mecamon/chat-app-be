package repositories_impl

import (
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestInitClusterMsgRepo(t *testing.T) {
	var i interface{}

	i = InitClusterMsgRepo(app, dbConn)

	if _, ok := i.(repositories.ClusterMsgRepo); !ok {
		t.Error("wrong type assertion")
	}
}

func TestGetClusterMsgRepo(t *testing.T) {
	var i interface{}

	i = GetClusterMsgRepo()

	if _, ok := i.(repositories.ClusterMsgRepo); !ok {
		t.Error("wrong type assertion")
	}
}

func TestClusterMsgRepo_Create(t *testing.T) {
	user := models.User{
		Name:      "user create cluster",
		Bio:       "user to create cluster",
		Email:     "usercreate@cluster.com",
		Password:  "validPass123",
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPass
	insertedUID, err := authTestRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	chat := models.GroupChat{
		Name:        "chat to create cluster",
		Description: "chat to create cluster",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}
	groupObjectID, err := primitive.ObjectIDFromHex(insertedGroupID)
	if err != nil {
		t.Error(err.Error())
	}

	var createTests = []struct {
		testName string
		cluster  models.ClusterOfMessages
		err      error
	}{
		{testName: "valid entry", cluster: models.ClusterOfMessages{
			BelongsToGroup: groupObjectID,
			CreatedAt:      time.Now().Unix(),
			UpdatedAt:      time.Now().Unix(),
			IsClosed:       false,
		}, err: nil},
	}

	for _, tt := range createTests {
		t.Log("TEST NAME:", tt.testName)

		_, err := clusterMsgRepo.Create(tt.cluster)
		if err != nil {
			t.Error(err.Error())
		}

		if tt.err == nil && err != nil {
			t.Errorf("was NOT expecting an error but got: %s", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Errorf("was expecting an error but got nothing")
		}
	}
}

func TestClusterMsgRepo_GetLatest(t *testing.T) {
	user := models.User{
		Name:      "User get latest",
		Bio:       "user get latest",
		Email:     "userget@latest.com",
		Password:  "validPass123",
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPass
	insertedUID, err := authTestRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	chat := models.GroupChat{
		Name:        "chat to get latest",
		Description: "chat to get latest",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	groupObjectID, err := primitive.ObjectIDFromHex(insertedGroupID)
	if err != nil {
		t.Error(err.Error())
	}

	cluster1 := models.ClusterOfMessages{
		BelongsToGroup: groupObjectID,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		IsClosed:       false,
	}
	cluster2 := models.ClusterOfMessages{
		BelongsToGroup: groupObjectID,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		IsClosed:       false,
	}

	_, err = clusterMsgRepo.Create(cluster1)
	if err != nil {
		t.Error(err.Error())
	}
	insertedClusterID2, err := clusterMsgRepo.Create(cluster2)
	if err != nil {
		t.Error(err.Error())
	}

	var getLatestTests = []struct {
		testName          string
		expectedClusterID string
		err               error
	}{
		{testName: "valid get", expectedClusterID: insertedClusterID2, err: nil},
	}

	for _, tt := range getLatestTests {
		t.Log("TEST NAME:", tt.testName)

		cluster, err := clusterMsgRepo.GetLatest()
		if tt.err == nil && err != nil {
			t.Errorf("it was NOT expecting an error but got: %s", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Error("it was expecting an error but got nothing")
		}

		clusterID := cluster.ID.Hex()
		if clusterID != tt.expectedClusterID {
			t.Error("clusterID is not equal to the last entry")
		}
	}
}

func TestClusterMsgRepo_Update(t *testing.T) {
	user := models.User{
		Name:      "user for update chat",
		Bio:       "user for update chat",
		Email:     "userfor@updatechat.com",
		Password:  "validPass123",
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPass
	insertedUID, err := authTestRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	chat := models.GroupChat{
		Name:        "chat for update",
		Description: "chat for update",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}
	groupObjectID, err := primitive.ObjectIDFromHex(insertedGroupID)
	if err != nil {
		t.Error(err.Error())
	}

	cluster := models.ClusterOfMessages{
		BelongsToGroup: groupObjectID,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		IsClosed:       false,
	}

	clusterID, err := clusterMsgRepo.Create(cluster)
	if err != nil {
		t.Error(err.Error())
	}

	messengerUID := "3eb3d668b31de5d588f42a6d"
	objectMsgUID, err := primitive.ObjectIDFromHex(messengerUID)
	if err != nil {
		t.Error(err.Error())
	}

	objectReceiverUID, err := primitive.ObjectIDFromHex(insertedUID)
	if err != nil {
		t.Error(err.Error())
	}

	message1 := models.Message{
		From:        objectMsgUID,
		To:          objectReceiverUID,
		TextContent: "test message content #1",
	}
	message2 := models.Message{
		From:        objectMsgUID,
		To:          objectReceiverUID,
		TextContent: "test message content #2",
	}

	var updateTests = []struct {
		testName  string
		clusterID string
		messages  []models.Message
		err       error
	}{
		{testName: "successful update", clusterID: clusterID, messages: []models.Message{
			message1,
			message2,
		}, err: nil},
	}

	for _, tt := range updateTests {
		t.Log("TEST NAME:", tt.testName)

		for _, m := range tt.messages {
			if err := clusterMsgRepo.Update(tt.clusterID, m); err != nil {
				t.Errorf("was NOT expecting an error but bot: %s", err.Error())
			}
		}
		clusterUpdated, err := clusterMsgRepo.GetLatest()
		if err != nil {
			t.Error(err.Error())
		}
		if len(clusterUpdated.Messages) != len(tt.messages) {
			t.Errorf("The length expected of the messages is %d, but got %d", len(tt.messages), len(clusterUpdated.Messages))
		}
	}
}

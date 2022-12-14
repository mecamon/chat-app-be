package repositories_impl

import (
	"errors"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestInitGroupChatRepo(t *testing.T) {
	var i interface{}
	i = InitGroupChatRepo(app, dbConn)

	if _, ok := i.(repositories.GroupChat); !ok {
		t.Error("error asserting type")
	}
}

func TestGetGroupChatRepo(t *testing.T) {
	var i interface{}
	i = GetGroupChatRepo()

	if _, ok := i.(repositories.GroupChat); !ok {
		t.Error("error asserting type")
	}
}

func TestGroupChatImpl_Create(t *testing.T) {
	var createTests = []struct {
		testName  string
		uid       string
		groupChat models.GroupChat
	}{
		{testName: "valid user", uid: "5eb3d668b31de5d588f42a7a", groupChat: models.GroupChat{
			Name:         "My group chat",
			Description:  "This is just a group",
			ImageURL:     "",
			Participants: nil,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}},
	}

	for _, tt := range createTests {
		t.Log("TEST NAME:", tt.testName)
		_, err := groupChatTestRepo.Create(tt.uid, tt.groupChat)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGroupChatImpl_Update(t *testing.T) {
	uid := "5eb3d668b31de5d588f42a7b"
	group := models.GroupChat{
		Name:         "Without update",
		Description:  "Group to update",
		ImageURL:     "",
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	insertedID, err := groupChatTestRepo.Create(uid, group)
	if err != nil {
		t.Error(err.Error())
	}

	ownerID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		t.Error(err.Error())
	}

	groupID, err := primitive.ObjectIDFromHex(insertedID)
	if err != nil {
		t.Error(err.Error())
	}

	validGroupU := models.GroupChat{
		ID:           groupID,
		Name:         "Updated",
		Description:  "Group already updated",
		ImageURL:     "",
		GroupOwner:   ownerID,
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	invalidGroupID := models.GroupChat{
		ID:         primitive.ObjectID{},
		GroupOwner: ownerID,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	invalidOwnerGroupU := models.GroupChat{
		ID:         primitive.ObjectID{},
		GroupOwner: ownerID,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	var updateTests = []struct {
		testName string
		groupU   models.GroupChat
		err      error
	}{
		{testName: "valid update", groupU: validGroupU, err: nil},
		{testName: "wrong groupID", groupU: invalidGroupID, err: errors.New("has error")},
		{testName: "wrong ownerID", groupU: invalidOwnerGroupU, err: errors.New("has error")},
	}

	for _, tt := range updateTests {
		t.Log("TEST NAME:", tt.testName)
		err = groupChatTestRepo.Update(tt.groupU)

		if tt.err == nil && err != nil {
			t.Error("error was NOT expected but got one:", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Error("expected error but did NOT get one")
		}
	}
}

func TestGroupChatImpl_Delete(t *testing.T) {
	uid := "5eb3d668b31de5d588f42a7c"
	wrongUid := "5eb3d668b31de5d588f42a7d"
	group := models.GroupChat{
		Name:         "Without update",
		Description:  "Group to update",
		ImageURL:     "",
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	insertedID, err := groupChatTestRepo.Create(uid, group)
	if err != nil {
		t.Error(err.Error())
	}

	var deleteTests = []struct {
		testName string
		ownerID  string
		groupID  string
		err      error
	}{
		{testName: "valid delete", ownerID: uid, groupID: insertedID, err: nil},
		{testName: "invalid delete", ownerID: wrongUid, groupID: "5eb3d668b31de5d588f42a7e", err: errors.New("has error")},
	}

	for _, tt := range deleteTests {
		t.Log("TEST NAME:", tt.testName)
		err := groupChatTestRepo.Delete(tt.ownerID, tt.groupID)

		if tt.err == nil && err != nil {
			t.Error("error was NOT expected but got one:", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Error("expected error but did NOT get one")
		}
	}
}

func TestGroupChatImpl_AddUserToChat(t *testing.T) {
	uid := "5eb3d668b31de5d588f42a8c"
	group := models.GroupChat{
		Name:         "Group to add user",
		Description:  "Group to add user",
		ImageURL:     "",
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(uid, group)
	if err != nil {
		t.Error(err.Error())
	}
	notExistingGroupID := "5eb3d668b31de5d588f42a7t"

	userToAdd := models.User{
		Name:      "User to add",
		Bio:       "I am an user to add to group",
		Email:     "userto@add.com",
		Password:  "validPass123",
		Phone:     8097656789,
		PhotoURL:  "",
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	_, err = authTestRepo.Register(userToAdd)
	if err != nil {
		t.Error(err.Error())
	}

	var addUserToChatTests = []struct {
		testName string
		groupID  string
		err      error
	}{
		{testName: "existing group ID", groupID: insertedGroupID, err: nil},
		{testName: "NOT existing group ID", groupID: notExistingGroupID, err: errors.New("has error")},
	}

	for _, tt := range addUserToChatTests {
		u, err := authTestRepo.FindByEmail(userToAdd.Email)
		if err != nil {
			t.Error(err.Error())
		}
		err = groupChatTestRepo.AddUserToChat(u, tt.groupID)

		if tt.err == nil && err != nil {
			t.Error("error was NOT expected but got:", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Error("error was expected but did not got one")
		}
	}
}

func TestGroupChatImpl_LoadAll(t *testing.T) {
	uid := "5eb3d668b31de5d588f42a7d"
	group := models.GroupChat{
		Name:         "Load all 1",
		Description:  "Group to load",
		ImageURL:     "",
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	groupID, err := groupChatTestRepo.Create(uid, group)
	if err != nil {
		t.Error(err.Error())
	}

	var skip int64 = 0
	var take int64 = 10

	participantData := models.User{
		Name:      "Participant",
		Bio:       "User to add as participant",
		Email:     "participant@add.com",
		Password:  "weooPass222",
		Phone:     8097651234,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	participantID, err := authTestRepo.Register(participantData)
	if err != nil {
		t.Error(err.Error())
	}
	participant, err := authTestRepo.FindByEmail(participantData.Email)
	if err != nil {
		t.Error(err.Error())
	}

	err = groupChatTestRepo.AddUserToChat(participant, groupID)
	if err != nil {
		t.Error(err.Error())
	}

	var loadAllTest = []struct {
		testName string
		uid      string
		filter   map[string]interface{}
		err      error
	}{
		{testName: "load owned groups only", uid: uid, filter: map[string]interface{}{
			"skip":  skip,
			"take":  take,
			"chats": "owned",
		}, err: nil},
		{testName: "load owned groups only", uid: participantID, filter: map[string]interface{}{
			"skip":  skip,
			"take":  take,
			"chats": "participating",
		}, err: nil},
	}

	for _, tt := range loadAllTest {
		groups, err := groupChatTestRepo.LoadAll(tt.uid, tt.filter)
		if tt.err == nil && err != nil {
			t.Error("did not expected an error but got one:", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Error("expected an error but did NOT get it")
		}

		if tt.filter["chats"] == "owned" && groups[0].GroupOwner.Hex() != tt.uid {
			t.Error("returned groups does NOT match the ownerID")
		}
		if tt.filter["chats"] == "participating" {
			gg := map[string]bool{}
			for i, g := range groups {
				gg[g.Name] = false
				for _, p := range groups[i].Participants {
					if p.Email == participant.Email {
						gg[g.Name] = true
						return
					}
				}
			}

			for _, v := range gg {
				if !v {
					t.Error("user is not in all groups returned")
				}
			}
		}
	}
}

func TestGroupChatImpl_AddImageURL(t *testing.T) {
	user := models.User{
		Name:      "add photo url user",
		Bio:       "add image url user",
		Email:     "addimageurl@user.com",
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
	notOwnerID := "5eb3d668b31de5d588f42a7d"

	chat := models.GroupChat{
		Name:        "group to add image",
		Description: "group to add image url",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var addImageURLTests = []struct {
		testName string
		uid      string
		groupID  string
		imageURL string
		err      error
	}{
		{
			testName: "not the owner of the group",
			uid:      notOwnerID,
			groupID:  insertedGroupID,
			imageURL: "https://image.com",
			err:      errors.New("has error"),
		},
		{
			testName: "valid update",
			uid:      insertedUID,
			groupID:  insertedGroupID,
			imageURL: "https://image.com",
			err:      nil,
		},
	}

	for _, tt := range addImageURLTests {
		err := groupChatTestRepo.AddImageURL(tt.uid, tt.groupID, tt.imageURL)
		if tt.err == nil && err != nil {
			t.Errorf("was not expecting an error but got: %s", err.Error())
		}
	}
}

func TestGroupChatImpl_RemoveImageURL(t *testing.T) {
	user := models.User{
		Name:      "remove photo url user",
		Bio:       "remove image url user",
		Email:     "removeimageurl@user.com",
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
	notOwnerID := "5eb3d668b31de5d588f42a7d"

	chat := models.GroupChat{
		Name:        "group to remove image",
		Description: "group to remove image url",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var removeImageURLTests = []struct {
		testName string
		uid      string
		groupID  string
		err      error
	}{
		{
			testName: "not the owner of the group",
			uid:      notOwnerID,
			groupID:  insertedGroupID,
			err:      errors.New("has error"),
		},
		{
			testName: "valid remove",
			uid:      insertedUID,
			groupID:  insertedGroupID,
			err:      nil,
		},
	}

	for _, tt := range removeImageURLTests {
		err := groupChatTestRepo.RemoveImageURL(tt.uid, tt.groupID)
		if tt.err == nil && err != nil {
			t.Errorf("was not expecting an error but got: %s", err.Error())
		}
	}
}

func TestGroupChatImpl_IsGroupOwner(t *testing.T) {
	user := models.User{
		Name:      "User to check owner",
		Bio:       "User created to check ownership",
		Email:     "user@isowner.com",
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
	notInsertedUID := "3eb3d668b31de5d588f42a5d"

	chat := models.GroupChat{
		Name:        "Group to check owner",
		Description: "Group to check ownership",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := groupChatTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var isGroupOwnerTests = []struct {
		testName       string
		uid            string
		groupID        string
		expectedResult bool
	}{
		{testName: "not the owner", uid: notInsertedUID, groupID: insertedGroupID, expectedResult: false},
		{testName: "it is the owner", uid: insertedUID, groupID: insertedGroupID, expectedResult: true},
	}

	for _, tt := range isGroupOwnerTests {
		t.Log("TEST NAME:", tt.testName)
		isOwner, _ := groupChatTestRepo.IsGroupOwner(tt.uid, tt.groupID)
		if isOwner != tt.expectedResult {
			t.Errorf("expected result is %v but got %v instead", tt.expectedResult, isOwner)
		}
	}
}

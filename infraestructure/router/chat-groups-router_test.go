package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	json_web_token "github.com/mecamon/chat-app-be/interface/json-web-token"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/interactors"
	"github.com/mecamon/chat-app-be/utils"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestChatGroupsController_Create(t *testing.T) {
	user := models.User{
		Name:      "Create chat ctrl",
		Bio:       "To create chat controller",
		Email:     "createchat@controller.com",
		Password:  "passwordValid123",
		Phone:     8091234567,
		PhotoURL:  "",
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

	var createTests = []struct {
		testName           string
		uid                string
		uEntry             models.GroupChat
		expectedStatusCode int
	}{
		{testName: "valid entry to create", uid: insertedUID, uEntry: models.GroupChat{
			Name:         "Create group",
			Description:  "The create ctrl",
			Participants: nil,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}, expectedStatusCode: http.StatusCreated},
		{testName: "invalid ID to create", uid: "1eb3d668b", uEntry: models.GroupChat{
			Name:         "Create group2",
			Description:  "The create ctrl2",
			Participants: nil,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}, expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range createTests {
		t.Log("TEST NAME:", tt.testName)
		token, err := json_web_token.Generate(tt.uid, "")
		if err != nil {
			t.Error(err.Error())
		}
		body, err := json.Marshal(tt.uEntry)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/chat_groups/", bytes.NewReader(body))
		req.Header.Set("Authorization", token)
		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is: %d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestChatGroupsController_Update(t *testing.T) {
	insertedUID, err := createUserForUpdate()
	if err != nil {
		t.Error(err.Error())
	}
	insertedChatID, err := createChatToUpdate(insertedUID)
	if err != nil {
		t.Error(err.Error())
	}

	var updateTests = []struct {
		testName           string
		uEntry             models.GroupChatDTO
		expectedStatusCode int
	}{
		{testName: "valid update", uEntry: models.GroupChatDTO{
			ID:          insertedChatID,
			Name:        "Name updated",
			Description: "Description updated",
			GroupOwner:  insertedUID,
		}, expectedStatusCode: http.StatusOK},
		{testName: "invalid ID type", uEntry: models.GroupChatDTO{
			ID:          "123123dasdas",
			Name:        "Name updated",
			Description: "Description updated",
			GroupOwner:  "1233asdasdasd",
		}, expectedStatusCode: http.StatusBadRequest},
		{testName: "invalid name and description", uEntry: models.GroupChatDTO{
			ID:          insertedChatID,
			Name:        utils.WordsGenerator(interactors.MinChatLengthName - 2),
			Description: utils.WordsGenerator(interactors.MaxChatLengthDes + 2),
			GroupOwner:  insertedUID,
		}, expectedStatusCode: http.StatusBadRequest},
		{
			testName: "not the user owner ID", uEntry: models.GroupChatDTO{
				ID:          insertedUID,
				Name:        "New name",
				Description: "New description",
				GroupOwner:  insertedChatID,
			}, expectedStatusCode: http.StatusForbidden},
	}

	for _, tt := range updateTests {
		t.Log("TEST NAME:", tt.testName)
		token, err := json_web_token.Generate(insertedUID, "")
		if err != nil {
			t.Error(err.Error())
		}

		body, err := json.Marshal(tt.uEntry)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/chat_groups/", bytes.NewReader(body))
		req.Header.Set("Authorization", token)
		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is: %d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
	}
}

func createUserForUpdate() (string, error) {
	user := models.User{
		Name:      "User update chat1",
		Bio:       "Update chat bio",
		Email:     "update@chatctrl.com",
		Password:  "passwordValid123",
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashPass
	insertedUID, err := authTestRepo.Register(user)
	if err != nil {
		return "", err
	}
	return insertedUID, nil
}

func createChatToUpdate(insertedUID string) (string, error) {
	chat := models.GroupChat{
		Name:         "Create group",
		Description:  "The create ctrl",
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	insertedChatID, err := chatGroupsTestRepo.Create(insertedUID, chat)
	if err != nil {
		return "", err
	}
	return insertedChatID, nil
}

func TestChatGroupsController_Delete(t *testing.T) {
	user := models.User{
		Name:      "Delete chat user",
		Bio:       "Delete chat bio",
		Email:     "deletechat@controller.com",
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
		Name:        "Chat to delete ctrl",
		Description: "Chat to delete ctrl",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := chatGroupsTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var deleteTests = []struct {
		testName           string
		groupID            string
		expectedStatusCode int
	}{
		{testName: "valid delete", groupID: insertedGroupID, expectedStatusCode: http.StatusOK},
		{testName: "invalid groupID to delete", groupID: "asdklbasdjbakjdb12312312", expectedStatusCode: http.StatusForbidden},
	}

	for _, tt := range deleteTests {
		t.Log("TEST NAME:", tt.testName)

		token, err := json_web_token.Generate(insertedUID, "")
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/chat_groups/%s", tt.groupID), nil)
		req.Header.Set("Authorization", token)
		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is: %d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestChatGroupsController_LoadAll(t *testing.T) {
	user := models.User{
		Name:      "User to load ctrl",
		Bio:       "User to load all controller",
		Email:     "userloadall@controller.com",
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
		Name:        "Chat to load all",
		Description: "Chat to load all",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	_, err = chatGroupsTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var loadAllTests = []struct {
		testName           string
		queryParams        string
		expectedStatusCode int
	}{
		{testName: "get all", queryParams: "?skip=0&take=8&chats=all", expectedStatusCode: http.StatusOK},
		{testName: "get my groups", queryParams: "?skip=0&take=8&chats=owned", expectedStatusCode: http.StatusOK},
		{testName: "get groups I'm participating", queryParams: "?skip=0&take=8&chats=participating", expectedStatusCode: http.StatusOK},
	}

	for _, tt := range loadAllTests {
		t.Log("TEST NAME:", tt.testName)
		token, err := json_web_token.Generate(insertedUID, user.Email)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chat_groups/load_chats%s", tt.queryParams), nil)
		req.Header.Set("Authorization", token)
		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is: %d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
		if tt.testName == "get groups I'm participating" {
			var groups []models.GroupChatDTO

			if err := json.NewDecoder(rr.Result().Body).Decode(&groups); err != nil {
				t.Error(err.Error())
			}
			if len(groups) != 1 {
				t.Error("expected 1 groups but got a nothing or more")
			}
		}
	}
}

func TestChatGroupsController_AddUserToChat(t *testing.T) {
	user := models.User{
		Name:      "user to add ctrl",
		Bio:       "user to add to controller",
		Email:     "usertoadd@controller.com",
		Password:  "validPass1234",
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPassword
	insertedUID, err := authTestRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	chat := models.GroupChat{
		Name:        "Chat to add user",
		Description: "Chat to add a user controller",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := chatGroupsTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var addToChatTests = []struct {
		testName           string
		groupID            string
		expectedStatusCode int
	}{
		{testName: "valid group id to add user", groupID: insertedGroupID, expectedStatusCode: http.StatusOK},
		{testName: "invalid group id to add user", groupID: "sadadad4234234", expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range addToChatTests {
		t.Log("TEST NAME:", tt.testName)

		token, err := json_web_token.Generate(insertedUID, user.Email)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/chat_groups/add_to_chat/%s", tt.groupID), nil)
		req.Header.Set("Authorization", token)
		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is %d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestChatGroupsController_AddImageURL(t *testing.T) {
	user := models.User{
		Name:      "add image to chat",
		Bio:       "User to test add image to chat",
		Email:     "addimage@tochat.com",
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
	notInsertedUID := "3eb3d668b31de5d588f42a6d"

	chat := models.GroupChat{
		Name:        "Chat to add image ctrl",
		Description: "Chat to add image ctrl",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := chatGroupsTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	wrongFile, err := os.OpenFile("../../fixtures/not-jpg.webp", os.O_RDONLY, 0755)
	if err != nil {
		t.Error(err.Error())
	}
	correctFile, err := os.OpenFile("../../fixtures/wildlife.jpg", os.O_RDONLY, 0755)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		wrongFile.Close()
		correctFile.Close()
	}()

	var addImageTests = []struct {
		testName           string
		uid                string
		groupID            string
		file               *os.File
		expectedStatusCode int
	}{
		{testName: "wrong file type", uid: insertedUID, groupID: insertedGroupID, file: wrongFile, expectedStatusCode: http.StatusBadRequest},
		{testName: "null field", uid: insertedUID, groupID: insertedGroupID, file: nil, expectedStatusCode: http.StatusBadRequest},
		{testName: "not the group owner", uid: notInsertedUID, groupID: insertedGroupID, file: correctFile, expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range addImageTests {
		t.Log("TEST NAME:", tt.testName)

		token, err := json_web_token.Generate(insertedUID, user.Email)
		if err != nil {
			t.Error(err.Error())
		}

		//adding file read from fixtures to multipart/form-data
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		if tt.file != nil {
			part, err := writer.CreateFormFile("file", tt.file.Name())
			if err != nil {
				t.Error(err.Error())
			}
			_, err = io.Copy(part, tt.file)
			if err != nil {
				t.Error(err.Error())
			}
		}

		writer.Close()

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/chat_groups/image/%s", tt.groupID), body)
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is:%d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestChatGroupsController_RemoveImageURL(t *testing.T) {
	user := models.User{
		Name:      "User to remove image ctrl",
		Bio:       "User to remove image ctrl",
		Email:     "usertoremove@imagectrl.com",
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
	notInsertedUID := "3eb3d668b31de5d588f42a6d"

	chat := models.GroupChat{
		Name:        "chat to remove image ctrl",
		Description: "chat to remove image ctrl",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedGroupID, err := chatGroupsTestRepo.Create(insertedUID, chat)
	if err != nil {
		t.Error(err.Error())
	}

	var removeImageURLTests = []struct {
		testName           string
		uid                string
		groupID            string
		expectedStatusCode int
	}{
		{testName: "not the owner", uid: notInsertedUID, groupID: insertedGroupID, expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range removeImageURLTests {
		t.Log("TEST NAME", tt.testName)

		token, err := json_web_token.Generate(tt.uid, user.Email)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/chat_groups/image/%s", tt.groupID), nil)
		req.Header.Set("Authorization", token)

		mainRouter.ServeHTTP(rr, req)

		if tt.expectedStatusCode != rr.Code {
			t.Errorf("expected status code is %d but got %d instead", tt.expectedStatusCode, rr.Code)
		}
	}

}

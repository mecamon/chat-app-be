//go:build !integration
// +build !integration

package interactors

import (
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
	"testing"
)

func TestEvalGroupInfo(t *testing.T) {
	var evalGroupInfoTests = []struct {
		testName             string
		groupInfo            models.GroupChat
		expectedNumberErrors int
	}{
		{testName: "name too short", groupInfo: models.GroupChat{Name: "d", Description: utils.WordsGenerator(minNameLength - 1)}, expectedNumberErrors: 1},
		{testName: "name too long", groupInfo: models.GroupChat{Name: utils.WordsGenerator(maxLengthName + 1), Description: "asodnasdandas"}, expectedNumberErrors: 1},
		{testName: "description too long", groupInfo: models.GroupChat{Name: "Done", Description: utils.WordsGenerator(maxLengthDescription + 1)}, expectedNumberErrors: 1},
		{testName: "valid group info", groupInfo: models.GroupChat{Name: "My group", Description: utils.WordsGenerator(60)}, expectedNumberErrors: 0},
	}

	for _, tt := range evalGroupInfoTests {
		t.Log("TEST NAME:", tt.testName)
		_, errors := EvalGroupInfo(tt.groupInfo)
		if len(errors) != tt.expectedNumberErrors {
			t.Errorf("errors expected are %d but got %d", tt.expectedNumberErrors, len(errors))
		}
	}
}

func TestCompleteGroupInfo(t *testing.T) {
	group := models.GroupChat{
		Name:        "My group",
		Description: "This is the group description",
	}
	completedGroup := CompleteGroupInfo(group)
	if completedGroup.CreatedAt == 0 || completedGroup.UpdatedAt == 0 {
		t.Error("CreatedAt or UpdatedAt information are 0")
	}
}

func TestGroupInfoToUpdate(t *testing.T) {
	uid := "5eb3d668b31de5d588f42a7a"
	group := models.GroupChatDTO{
		ID:           "5eb3d668b31de5d588f42a7a",
		Name:         "My chat",
		Description:  "This is a chat",
		ImageURL:     "",
		GroupOwner:   uid,
		Participants: nil,
	}

	groupU, err := GroupInfoToUpdate(uid, group)
	if err != nil {
		t.Error(err.Error())
	}

	if groupU.ID.Hex() == "" {
		t.Error("could not get the chat ID")
	}
	if groupU.GroupOwner.Hex() == "" {
		t.Error("could not get the chat owner ID")
	}
	if groupU.UpdatedAt == 0 {
		t.Error("UpdatedAt has not been set")
	}
}

//go:build !integration
// +build !integration

package presenters

import (
	"github.com/mecamon/chat-app-be/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestFormatGroups(t *testing.T) {
	groupA := models.GroupChat{
		ID:           primitive.ObjectID{},
		Name:         "Group a",
		Description:  "random group",
		ImageURL:     "",
		GroupOwner:   primitive.ObjectID{},
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	groupB := models.GroupChat{
		ID:           primitive.ObjectID{},
		Name:         "Group B",
		Description:  "random group",
		ImageURL:     "",
		GroupOwner:   primitive.ObjectID{},
		Participants: nil,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	groups := []models.GroupChat{
		groupA,
		groupB,
	}

	var i interface{}

	i = FormatGroups(groups)

	if _, ok := i.([]models.GroupChatDTO); !ok {
		t.Errorf("wrong type returned")
	}
}

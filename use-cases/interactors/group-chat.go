package interactors

import (
	"github.com/mecamon/chat-app-be/models"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	MinChatLengthName int = 2
	MaxChatLengthName int = 30
	MaxChatLengthDes  int = 100
)

func EvalGroupInfo(group models.GroupChat) (bool, []*cErrors.Custom) {
	var errSlices []*cErrors.Custom

	if len(group.Name) > MaxChatLengthName {
		errSlices = append(errSlices, &cErrors.Custom{
			Property:     "name",
			MessageID:    "NameTooLong",
			TemplateData: map[string]interface{}{"Count": MaxChatLengthName},
		})
	}

	if len(group.Name) < MinChatLengthName {
		errSlices = append(errSlices, &cErrors.Custom{
			Property:     "name",
			MessageID:    "NameTooShort",
			TemplateData: map[string]interface{}{"Count": minNameLength},
		})
	}

	if len(group.Description) > MaxChatLengthDes {
		errSlices = append(errSlices, &cErrors.Custom{
			Property:     "description",
			MessageID:    "DescriptionTooLong",
			TemplateData: map[string]interface{}{"Count": MaxChatLengthDes},
		})
	}

	return len(errSlices) == 0, errSlices
}

func CompleteGroupInfo(group models.GroupChat) models.GroupChat {
	group.CreatedAt = time.Now().Unix()
	group.UpdatedAt = time.Now().Unix()
	return group
}

func GroupInfoToUpdate(uid string, group models.GroupChatDTO) (models.GroupChat, error) {
	var groupU models.GroupChat

	OwnerID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return groupU, err
	}
	ID, err := primitive.ObjectIDFromHex(group.ID)
	if err != nil {
		return groupU, err
	}

	groupU = models.GroupChat{
		ID:           ID,
		Name:         group.Name,
		Description:  group.Description,
		ImageURL:     group.ImageURL,
		GroupOwner:   OwnerID,
		Participants: nil,
		CreatedAt:    0,
		UpdatedAt:    time.Now().Unix(),
	}
	return groupU, nil
}

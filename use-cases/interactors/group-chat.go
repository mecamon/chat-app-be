package interactors

import (
	"github.com/mecamon/chat-app-be/models"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	minLengthName        int = 2
	maxLengthName        int = 30
	maxLengthDescription int = 100
)

func EvalGroupInfo(group models.GroupChat) (bool, []*cErrors.Custom) {
	var errSlices []*cErrors.Custom

	if len(group.Name) > maxLengthName {
		errSlices = append(errSlices, &cErrors.Custom{
			Property:     "name",
			MessageID:    "NameTooLong",
			TemplateData: map[string]interface{}{"Count": maxLengthName},
		})
	}

	if len(group.Name) < minLengthName {
		errSlices = append(errSlices, &cErrors.Custom{
			Property:     "name",
			MessageID:    "NameTooShort",
			TemplateData: map[string]interface{}{"Count": minNameLength},
		})
	}

	if len(group.Description) > maxLengthDescription {
		errSlices = append(errSlices, &cErrors.Custom{
			Property:     "description",
			MessageID:    "DescriptionTooLong",
			TemplateData: map[string]interface{}{"Count": maxLengthDescription},
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

	ID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return groupU, err
	}
	OwnerID, err := primitive.ObjectIDFromHex(group.GroupOwner)
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

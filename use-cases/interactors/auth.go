package interactors

import (
	"fmt"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
	"time"

	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
)

const (
	minNameLength  = 2
	minPhoneLength = 9
)

func EvalRegistryEntry(uEntry models.User) (bool, []*cErrors.Custom) {
	var errors []*cErrors.Custom

	if len(uEntry.Name) < minNameLength {
		errors = append(errors, &cErrors.Custom{
			Property:     "name",
			MessageID:    "NameTooShort",
			TemplateData: map[string]interface{}{"Count": minNameLength},
		})
	}

	if hasValidEmail := utils.HasValidEmail(uEntry.Email); !hasValidEmail {
		errors = append(errors, &cErrors.Custom{
			Property:     "email",
			MessageID:    "InvalidEmail",
			TemplateData: nil,
		})
	}

	if hasValidPassword := utils.HasValidPass(uEntry.Password); !hasValidPassword {
		errors = append(errors, &cErrors.Custom{
			Property:     "password",
			MessageID:    "InvalidPassword",
			TemplateData: nil,
		})
	}

	if len(fmt.Sprintf("%d", uEntry.Phone)) < minPhoneLength {
		errors = append(errors, &cErrors.Custom{
			Property:     "phone",
			MessageID:    "PhoneTooShort",
			TemplateData: map[string]interface{}{"Count": minPhoneLength},
		})
	}

	return len(errors) == 0, errors
}

func CompleteRegEntry(uEntry models.User) models.User {
	uEntry.IsActive = true
	uEntry.CreatedAt = time.Now().Unix()
	uEntry.UpdatedAt = time.Now().Unix()
	return uEntry
}

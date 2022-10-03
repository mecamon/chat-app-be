package services

import (
	"fmt"
	"github.com/mecamon/chat-app-be/config"
	json_web_token "github.com/mecamon/chat-app-be/interface/json-web-token"
	"github.com/mecamon/chat-app-be/models"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
	"github.com/mecamon/chat-app-be/use-cases/interactors"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"
)

var auth *Auth

type Auth struct {
	app      *config.App
	authRepo repositories.AuthRepo
}

func InitAuth(app *config.App, repo repositories.AuthRepo) *Auth {
	auth := &Auth{
		app:      app,
		authRepo: repo,
	}
	return auth
}

func GetAuth() *Auth {
	return auth
}

func (a *Auth) Register(uEntry models.User) (string, []*cErrors.Custom) {
	_, errSlice := interactors.EvalRegistryEntry(uEntry)
	if errSlice != nil {
		return "", errSlice
	}
	completedU := interactors.CompleteRegEntry(uEntry)

	insertedID, err := a.authRepo.Register(completedU)
	if err != nil {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "email",
			MessageID:    "EmailAddressTaken",
			TemplateData: nil,
		})
		return "", errSlice
	}
	return insertedID, errSlice
}

func (a *Auth) Login(email, password string) (string, []*cErrors.Custom) {
	var errSlice []*cErrors.Custom

	ID, err := a.authRepo.Login(email, password)
	if err != nil {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "email",
			MessageID:    "InvalidEmailOrPassword",
			TemplateData: nil,
		})
		return "", errSlice
	}
	return ID, errSlice
}

func (a *Auth) SendRecoverPassLink(email string) (string, []*cErrors.Custom) {
	var errSlice []*cErrors.Custom

	hasValidEmail := utils.HasValidEmail(email)
	if !hasValidEmail {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "email",
			MessageID:    "InvalidEmail",
			TemplateData: nil,
		})
		return "", errSlice
	}

	user, err := a.authRepo.FindByEmail(email)
	if err != nil {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "email",
			MessageID:    "EmailDoesNotExist",
			TemplateData: nil,
		})
		return "", errSlice
	}

	ID := user.ID.Hex()
	token, err := json_web_token.Generate(ID, "")
	if err != nil {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "email",
			MessageID:    "ServerError",
			TemplateData: nil,
		})
		return "", errSlice
	}
	recoverLinkWithToken := fmt.Sprintf("%s?t=%s", a.app.RecoverHostAndPath, token)
	return recoverLinkWithToken, errSlice
}

func (a *Auth) ChangePassword(ID, newPassword string) []*cErrors.Custom {
	var errSlice []*cErrors.Custom

	if hasValidPass := utils.HasValidPass(newPassword); !hasValidPass {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "password",
			MessageID:    "InvalidPassword",
			TemplateData: nil,
		})
		return errSlice
	}
	if err := a.authRepo.ChangePassword(ID, newPassword); err != nil {
		errSlice = append(errSlice, &cErrors.Custom{
			Property:     "password",
			MessageID:    "ErrorChangingPass",
			TemplateData: nil,
		})
		return errSlice
	}

	return errSlice
}

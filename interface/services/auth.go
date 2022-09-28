package services

import (
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/models"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
	"github.com/mecamon/chat-app-be/use-cases/interactors"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
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

package services

import (
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
)

type Email struct {
	app      *config.App
	authRepo repositories.AuthRepo
}

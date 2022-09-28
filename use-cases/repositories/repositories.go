package repositories

import "github.com/mecamon/chat-app-be/models"

type AuthRepo interface {
	Register(models.User) (string, error)
}

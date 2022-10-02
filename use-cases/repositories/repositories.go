package repositories

import "github.com/mecamon/chat-app-be/models"

type AuthRepo interface {
	Register(models.User) (string, error)
	Login(email, password string) (string, error)
	FindByEmail(email string) (models.User, error)
	ChangePassword(id, newPassword string) error
}

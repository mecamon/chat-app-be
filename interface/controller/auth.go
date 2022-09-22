package controller

import (
	"net/http"

	"github.com/mecamon/chat-app-be/config"
)

type AuthController struct {
	app *config.App
}

func GetAuthController() *AuthController {
	return &AuthController{
		app: config.GetConfig(),
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("register endpoint"))
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("login endpoint"))
}
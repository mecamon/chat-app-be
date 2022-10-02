package controller

import (
	"encoding/json"
	appi18n "github.com/mecamon/chat-app-be/i18n"
	json_web_token "github.com/mecamon/chat-app-be/interface/json-web-token"
	"github.com/mecamon/chat-app-be/interface/services"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/presenters"
	"github.com/mecamon/chat-app-be/utils"
	"net/http"

	"github.com/mecamon/chat-app-be/config"
)

type AuthController struct {
	app          *config.App
	mLocales     *appi18n.MultiLocales
	authService  *services.Auth
	emailService *services.Email
}

var auth *AuthController

func InitAuthController(app *config.App, loc *appi18n.MultiLocales, authServ *services.Auth) *AuthController {
	auth = &AuthController{
		app:         app,
		mLocales:    loc,
		authService: authServ,
	}
	return auth
}

func GetAuthController() *AuthController {
	return auth
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)

	var uEntry models.User
	if err := json.NewDecoder(r.Body).Decode(&uEntry); err != nil {
		panic(w)
	}

	insertedID, errSlice := c.authService.Register(uEntry)
	if len(errSlice) > 0 {
		errMessages := presenters.ErrMessages(locales, errSlice)
		body, err := json.Marshal(errMessages)
		if err != nil {
			panic(w)
		}
		utils.Response(w, http.StatusBadRequest, body)
		return
	}

	token, err := json_web_token.Generate(insertedID, "")
	if err != nil {
		panic(w)
	}
	var regSuccess = struct {
		Token string `json:"token"`
	}{Token: token}

	body, err := json.Marshal(regSuccess)
	if err != nil {
		panic(w)
	}
	utils.Response(w, http.StatusCreated, body)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)

	var uEntry struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&uEntry); err != nil {
		panic(w)
	}

	ID, errColl := c.authService.Login(uEntry.Email, uEntry.Password)
	if len(errColl) != 0 {
		errMessages := presenters.ErrMessages(locales, errColl)
		body, err := json.Marshal(errMessages)
		if err != nil {
			panic(w)
		}
		utils.Response(w, http.StatusBadRequest, body)
		return
	}

	token, err := json_web_token.Generate(ID, "")
	if err != nil {
		panic(w)
	}
	regSuccess := struct {
		Token string `json:"token"`
	}{Token: token}

	body, err := json.Marshal(regSuccess)
	if err != nil {
		panic(w)
	}
	utils.Response(w, http.StatusOK, body)
}

func (c *AuthController) SendRecoveryLink(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)

	var uEntry struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&uEntry); err != nil {
		panic(w)
	}

	link, errSlice := c.authService.SendRecoverPassLink(uEntry.Email)
	if len(errSlice) != 0 {
		for _, ee := range errSlice {
			if ee.MessageID == "ServerError" {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}
			if ee.MessageID == "EmailDoesNotExist" {
				utils.Response(w, http.StatusNotFound, nil)
				return
			}
		}
		errMessages := presenters.ErrMessages(locales, errSlice)
		body, err := json.Marshal(errMessages)
		if err != nil {
			panic(w)
		}
		utils.Response(w, http.StatusBadRequest, body)
		return
	}

	resSuccess := struct {
		Link string `json:"link"`
	}{Link: link}

	body, err := json.Marshal(resSuccess)
	if err != nil {
		panic(w)
	}
	utils.Response(w, http.StatusOK, body)
}

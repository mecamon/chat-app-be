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
	app         *config.App
	mLocales    *appi18n.MultiLocales
	authService *services.Auth
}

var auth *AuthController

func InitAuthController(
	app *config.App,
	loc *appi18n.MultiLocales,
	authServ *services.Auth) *AuthController {
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
		errMsg := locales.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
	}

	insertedID, errSlice := c.authService.Register(uEntry)
	if len(errSlice) > 0 {
		errMessages := presenters.ErrMessages(locales, errSlice)
		if err := utils.JSONResponse(w, http.StatusBadRequest, errMessages); err != nil {
			panic(err)
		}
		return
	}

	token, err := json_web_token.Generate(insertedID, "")
	if err != nil {
		panic(w)
	}
	var regSuccess = struct {
		Token string `json:"token"`
	}{Token: token}

	if err := utils.JSONResponse(w, http.StatusCreated, regSuccess); err != nil {
		panic(err)
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)

	var uEntry struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&uEntry); err != nil {
		errMsg := locales.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
	}

	ID, errColl := c.authService.Login(uEntry.Email, uEntry.Password)
	if len(errColl) != 0 {
		errMessages := presenters.ErrMessages(locales, errColl)
		if err := utils.JSONResponse(w, http.StatusBadRequest, errMessages); err != nil {
			panic(err)
		}
		return
	}

	token, err := json_web_token.Generate(ID, "")
	if err != nil {
		panic(w)
	}
	regSuccess := struct {
		Token string `json:"token"`
	}{Token: token}

	if err := utils.JSONResponse(w, http.StatusOK, regSuccess); err != nil {
		panic(err)
	}
}

func (c *AuthController) SendRecoveryLink(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)

	var uEntry struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&uEntry); err != nil {
		errMsg := locales.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
	}

	_, errSlice := c.authService.SendRecoverPassLink(uEntry.Email)
	if len(errSlice) != 0 {
		for _, ee := range errSlice {
			if ee.MessageID == "ServerError" {
				_ = utils.JSONResponse(w, http.StatusInternalServerError, nil)
				return
			}
			if ee.MessageID == "EmailDoesNotExist" {
				_ = utils.JSONResponse(w, http.StatusNotFound, nil)
				return
			}
		}
		errMessages := presenters.ErrMessages(locales, errSlice)
		if err := utils.JSONResponse(w, http.StatusBadRequest, errMessages); err != nil {
			panic(err)
		}
		return
	}

	//TODO: send email with the link the user

	//resSuccess := struct {
	//	Link string `json:"link"`
	//}{Link: link}

	if err := utils.JSONResponse(w, http.StatusOK, nil); err != nil {
		panic(err)
	}
}

func (c *AuthController) ChangePass(w http.ResponseWriter, r *http.Request) {
	customClaims, err := json_web_token.Validate(r.Header.Get("Authorization"))
	if err != nil {
		if err := utils.JSONResponse(w, http.StatusUnauthorized, nil); err != nil {
			panic(err)
		}
		return
	}

	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)
	ID := customClaims.ID

	body := struct {
		NewPassword string `json:"newPassword"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errMsg := locales.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
	}

	errSlice := c.authService.ChangePassword(ID, body.NewPassword)
	if len(errSlice) != 0 {
		errMessages := presenters.ErrMessages(locales, errSlice)
		if err := utils.JSONResponse(w, http.StatusBadRequest, errMessages); err != nil {
			panic(err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, nil); err != nil {
		panic(err)
	}
}

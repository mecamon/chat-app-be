package controller

import (
	"encoding/json"
	"fmt"
	"github.com/mecamon/chat-app-be/config"
	appi18n "github.com/mecamon/chat-app-be/i18n"
	json_web_token "github.com/mecamon/chat-app-be/interface/json-web-token"
	"github.com/mecamon/chat-app-be/interface/services"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/interactors"
	"github.com/mecamon/chat-app-be/use-cases/presenters"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"
	"log"
	"net/http"
	"strconv"
)

const (
	maxFileSize int64 = 5242880
)

var fileAcceptedContentTypes = []string{"image/jpg", "image/jpeg", "image/png"}

type AuthController struct {
	app      *config.App
	mLocales *appi18n.MultiLocales
	authRepo repositories.AuthRepo
}

var auth *AuthController

func InitAuthController(app *config.App, loc *appi18n.MultiLocales, authRepo repositories.AuthRepo) *AuthController {
	auth = &AuthController{
		app:      app,
		mLocales: loc,
		authRepo: authRepo,
	}
	return auth
}

func GetAuthController() *AuthController {
	return auth
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := c.mLocales.GetSpeLocales(lang)

	if err := r.ParseMultipartForm(128); err != nil {
		errMsg := locales.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	uEntry := models.User{
		Name:     r.Form.Get("name"),
		Bio:      r.Form.Get("bio"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	file, fileHeader, _ := r.FormFile("file")

	phoneStr := r.Form.Get("phone")
	phone, err := strconv.ParseInt(phoneStr, 10, 0)
	if err != nil {
		errMsg := locales.GetMsg("PhoneWrongFormat", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}
	uEntry.Phone = phone

	_, errSlice := interactors.EvalRegistryEntry(uEntry)
	if len(errSlice) != 0 {
		errMessages := presenters.ErrMessages(locales, errSlice)
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	} else if len(errSlice) != 0 && file != nil {
		fileInfo := models.FileInfo{
			Size:        fileHeader.Size,
			ContentType: fileHeader.Header.Get("Content-Type"),
		}

		hasAValidFile := interactors.ValidFile(fileInfo, maxFileSize, fileAcceptedContentTypes...)
		if !hasAValidFile {
			errMsg := locales.GetMsg("WrongFileType", map[string]interface{}{
				"Types": fileAcceptedContentTypes,
				"Size":  maxFileSize,
			})
			errMessages := []string{errMsg}
			_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
			return
		}

		storage, err := services.GetStorage()
		if err != nil {
			log.Println(err.Error())
		}
		photoURL, err := storage.UploadImage(file, uEntry.Email)
		if err != nil {
			log.Println(err.Error())
		}
		uEntry.PhotoURL = photoURL
	}

	completedU := interactors.CompleteRegEntry(uEntry)
	insertedID, err := c.authRepo.Register(completedU)
	if err != nil {
		errMsg := locales.GetMsg("EmailAddressTaken", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusConflict, errMessages)
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

	ID, err := c.authRepo.Login(uEntry.Email, uEntry.Password)
	if err != nil {
		errMsg := locales.GetMsg("InvalidEmailOrPassword", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
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

	hasValidEmail := utils.HasValidEmail(uEntry.Email)
	if !hasValidEmail {
		errMsg := locales.GetMsg("InvalidEmail", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	user, err := c.authRepo.FindByEmail(uEntry.Email)
	if err != nil {
		errMsg := locales.GetMsg("EmailDoesNotExist", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusNotFound, errMessages)
	}

	ID := user.ID.Hex()
	token, err := json_web_token.Generate(ID, "")
	if err != nil {
		panic(w)
	}

	//the link formatted
	_ = fmt.Sprintf("%s?t=%s", c.app.RecoverHostAndPath, token)

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

	uEntry := struct {
		NewPassword string `json:"newPassword"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&uEntry); err != nil {
		errMsg := locales.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
	}

	if hasValidPass := utils.HasValidPass(uEntry.NewPassword); !hasValidPass {
		errMsg := locales.GetMsg("InvalidPassword", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	if err := c.authRepo.ChangePassword(ID, uEntry.NewPassword); err != nil {
		errMsg := locales.GetMsg("ErrorChangingPass", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, nil); err != nil {
		panic(err)
	}
}

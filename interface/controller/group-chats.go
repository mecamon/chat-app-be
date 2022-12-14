package controller

import (
	"encoding/json"
	"github.com/mecamon/chat-app-be/config"
	appi18n "github.com/mecamon/chat-app-be/i18n"
	repositories_impl "github.com/mecamon/chat-app-be/interface/repositories"
	"github.com/mecamon/chat-app-be/interface/services"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/interactors"
	"github.com/mecamon/chat-app-be/use-cases/presenters"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"github.com/mecamon/chat-app-be/utils"
	"net/http"
	"strconv"
)

var groupChats GroupChats

const maxGroupImageSize int64 = 5242880

type GroupChats struct {
	app           *config.App
	mLocales      *appi18n.MultiLocales
	autRepo       repositories.AuthRepo
	groupChatRepo repositories.GroupChat
}

func InitGroupChats() GroupChats {
	groupChats = GroupChats{
		app:           config.GetConfig(),
		mLocales:      appi18n.GetMultiLocales(),
		autRepo:       repositories_impl.GetAuthRepo(),
		groupChatRepo: repositories_impl.GetGroupChatRepo(),
	}
	return groupChats
}

func GetGroupChats() GroupChats {
	return groupChats
}

func (c *GroupChats) Create(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)

	var group models.GroupChat

	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		errMsg := loc.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	hasValidInfo, errSlice := interactors.EvalGroupInfo(group)
	if !hasValidInfo {
		errMessages := presenters.ErrMessages(loc, errSlice)
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}
	completedGroup := interactors.CompleteGroupInfo(group)

	insertedGroupID, err := c.groupChatRepo.Create(ID, completedGroup)
	if err != nil {
		errMsg := loc.GetMsg("ErrorCreatingGroup", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	success := struct {
		InsertedID string `json:"insertedID"`
	}{InsertedID: insertedGroupID}

	_ = utils.JSONResponse(w, http.StatusCreated, success)
}

func (c *GroupChats) Update(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)

	var group models.GroupChatDTO

	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		errMsg := loc.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	groupU, err := interactors.GroupInfoToUpdate(ID, group)
	if err != nil {
		errMsg := loc.GetMsg("IDInvalidType", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	isValid, errSlice := interactors.EvalGroupInfo(groupU)
	if !isValid {
		errMessages := presenters.ErrMessages(loc, errSlice)
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	if err := c.groupChatRepo.Update(groupU); err != nil {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusForbidden, errMessages)
		return
	}

	_ = utils.JSONResponse(w, http.StatusOK, nil)
}

func (c *GroupChats) Delete(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)
	groupID := utils.GetRouteParam(r.URL.Path)

	if err := c.groupChatRepo.Delete(ID, groupID); err != nil {
		errMsg := loc.GetMsg("NothingWasDeleted", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusForbidden, errMessages)
		return
	}

	_ = utils.JSONResponse(w, http.StatusOK, nil)
}

func (c *GroupChats) LoadAll(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)
	qParams := r.URL.Query()

	skip, err := strconv.ParseInt(qParams["skip"][0], 10, 64)
	if err != nil {
		skip = 0
	}
	take, err := strconv.ParseInt(qParams["take"][0], 10, 64)
	if err != nil {
		take = 10
	}
	chats := qParams["chats"][0]

	filters := map[string]interface{}{
		"skip":  skip,
		"take":  take,
		"chats": chats,
	}

	groups, err := c.groupChatRepo.LoadAll(ID, filters)
	if err != nil {
		errMsg := loc.GetMsg("ErrorGettingChatGroups", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusServiceUnavailable, errMessages)
		return
	}

	groupsF := presenters.FormatGroups(groups)

	if groupsF == nil || len(groupsF) == 0 {
		_ = utils.JSONResponse(w, http.StatusNoContent, groupsF)
		return
	}

	_ = utils.JSONResponse(w, http.StatusOK, groupsF)
}

func (c *GroupChats) AddUserToChat(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)
	groupID := utils.GetRouteParam(r.URL.Path)

	user, err := c.autRepo.FindByID(ID)
	if err != nil {
		_ = utils.JSONResponse(w, http.StatusForbidden, nil)
		return
	}

	if err := c.groupChatRepo.AddUserToChat(user, groupID); err != nil {
		errMsg := loc.GetMsg("ErrorAddingUserToChat", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	_ = utils.JSONResponse(w, http.StatusOK, nil)
}

func (c *GroupChats) AddImageURL(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)
	groupID := utils.GetRouteParam(r.URL.Path)

	if err := r.ParseMultipartForm(128); err != nil {
		errMsg := loc.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		errMsg := loc.GetMsg("ErrorParsingBody", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}
	defer file.Close()

	if file == nil {
		errMsg := loc.GetMsg("RequiredField", map[string]interface{}{"RequiredField": "file"})
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	fileInfo := models.FileInfo{
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	contentTypesAccepted := []string{
		"image/jpg",
		"image/jpeg",
		"image/png",
		"image/JPG",
		"image/JPEG",
		"image/PNG",
	}

	isValid := interactors.ValidFile(fileInfo, maxGroupImageSize, contentTypesAccepted...)
	if !isValid {
		errMsg := loc.GetMsg("WrongFileType", map[string]interface{}{
			"Types": contentTypesAccepted,
			"Size":  maxGroupImageSize,
		})
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	isOwner, err := c.groupChatRepo.IsGroupOwner(ID, groupID)
	if err != nil {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}
	if !isOwner {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	storage, err := services.GetStorage()
	if err != nil {
		_ = utils.JSONResponse(w, http.StatusServiceUnavailable, nil)
		return
	}

	imageURL, err := storage.UploadImage(file, groupID)
	if err != nil {
		errMsg := loc.GetMsg("ErrorAddingFile", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusServiceUnavailable, errMessages)
		return
	}

	if err := c.groupChatRepo.AddImageURL(ID, groupID, imageURL); err != nil {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	_ = utils.JSONResponse(w, http.StatusOK, nil)
}

func (c *GroupChats) RemoveImageURL(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	loc := c.mLocales.GetSpeLocales(lang)
	ID := r.Context().Value("ID").(string)
	groupID := utils.GetRouteParam(r.URL.Path)

	isOwner, err := c.groupChatRepo.IsGroupOwner(ID, groupID)
	if err != nil {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}
	if !isOwner {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	storage, err := services.GetStorage()
	if err != nil {
		errMsg := loc.GetMsg("ErrorRemovingFile", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusServiceUnavailable, errMessages)
		return
	}

	_, err = storage.DeleteImage(groupID)
	if err != nil {
		errMsg := loc.GetMsg("ErrorRemovingFile", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusServiceUnavailable, errMessages)
		return
	}

	if err := c.groupChatRepo.RemoveImageURL(ID, groupID); err != nil {
		errMsg := loc.GetMsg("NothingWasUpdated", nil)
		errMessages := []string{errMsg}
		_ = utils.JSONResponse(w, http.StatusBadRequest, errMessages)
		return
	}

	_ = utils.JSONResponse(w, http.StatusOK, nil)
}

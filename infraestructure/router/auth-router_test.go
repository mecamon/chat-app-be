//go:build integration
// +build integration

package router

import (
	"bytes"
	"encoding/json"
	json_web_token "github.com/mecamon/chat-app-be/interface/json-web-token"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthController_Register(t *testing.T) {
	var registerTests = []struct {
		testName       string
		uEntry         models.User
		expectedCode   int
		expectedErrors int
	}{
		{testName: "valid user credentials", uEntry: models.User{
			Name:     "Juan Lopez",
			Bio:      "Yo soy Juan",
			Email:    "juan@mail.com",
			Password: "JuanMan123456",
			Phone:    8097654312,
		}, expectedCode: http.StatusCreated, expectedErrors: 0},
		{testName: "email in use", uEntry: models.User{
			Name:     "Juan Lopez",
			Bio:      "Yo soy Juan",
			Email:    "juan@mail.com",
			Password: "JuanMan123456",
			Phone:    8097654312,
		}, expectedCode: http.StatusBadRequest, expectedErrors: 1},
		{testName: "no name", uEntry: models.User{
			Name:     "",
			Bio:      "Yo soy Juan",
			Email:    "juan@mail.com",
			Password: "JuanMan123456",
			Phone:    8097654312,
		}, expectedCode: http.StatusBadRequest, expectedErrors: 1},
		{testName: "wrong formatted email and password", uEntry: models.User{
			Name:     "Popa",
			Bio:      "Yo soy Juan",
			Email:    "juanmail.com",
			Password: "12345678676556",
			Phone:    8097654312,
		}, expectedCode: http.StatusBadRequest, expectedErrors: 2},
	}

	for _, tt := range registerTests {
		t.Log(tt.testName)
		body, _ := json.Marshal(&tt.uEntry)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		mainRouter.ServeHTTP(rr, req)

		if rr.Code != tt.expectedCode {
			t.Errorf("expected status code %d, but got %d", tt.expectedCode, rr.Code)
		}

		if rr.Code == http.StatusCreated {
			resBody := struct {
				Token string `json:"token"`
			}{}
			if err := json.NewDecoder(rr.Body).Decode(&resBody); err == nil {
				if resBody.Token == "" {
					t.Error("expected a token but got an empty string")
				}
			}
		}

		if tt.expectedCode == http.StatusBadRequest || tt.expectedCode == http.StatusConflict {
			var errResponse []string
			if err := json.NewDecoder(rr.Body).Decode(&errResponse); err != nil {
				t.Error(err.Error())
			}
			if tt.expectedErrors != len(errResponse) {
				t.Errorf("expected errors are %d, but got %d", tt.expectedErrors, len(errResponse))
			}
		}
	}
}

func TestAuthController_Login(t *testing.T) {
	password := "validPass1234"
	user := models.User{
		Name:      "Login cont user",
		Bio:       "This is the login controller user",
		Email:     "loginctrl@mail.com",
		Password:  password,
		Phone:     123455677655,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashPass, err := utils.GenerateHash(password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashPass

	_, err = authTestRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	var loginTests = []struct {
		testName, email, password string
		statusCode                int
	}{
		{testName: "successful login", email: user.Email, password: password, statusCode: http.StatusOK},
		{testName: "wrong email", email: "dasdadddafa", password: password, statusCode: http.StatusBadRequest},
		{testName: "wrong password", email: user.Email, password: "dasdaddf23", statusCode: http.StatusBadRequest},
	}

	for _, tt := range loginTests {
		t.Log(tt.testName)

		uEntry := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    tt.email,
			Password: tt.password,
		}

		body, _ := json.Marshal(uEntry)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		mainRouter.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK {
			resBody := struct {
				Token string `json:"token"`
			}{}
			if err := json.NewDecoder(rr.Body).Decode(&resBody); err == nil {
				if resBody.Token == "" {
					t.Error("expected a token but got an empty string")
				}
			}
		}

		if rr.Code != tt.statusCode {
			t.Errorf("expected statusCode is %d, but got %d", tt.statusCode, rr.Code)
		}
	}
}

//func TestAuthController_SendRecoveryLink(t *testing.T) {
//	password := "Password1234"
//	user := models.User{
//		Name:      "Send Recover Ctrl",
//		Bio:       "This is the send recover ctrl",
//		Email:     "sendrecover@controller.com",
//		Password:  password,
//		Phone:     809123456789,
//		PhotoURL:  "",
//		IsActive:  true,
//		CreatedAt: time.Now().Unix(),
//		UpdatedAt: time.Now().Unix(),
//	}
//	hashPass, err := utils.GenerateHash(password)
//	if err != nil {
//		t.Error(err.Error())
//	}
//	user.Password = hashPass
//
//	_, err = authTestRepo.Register(user)
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	var sendRecoverTests = []struct {
//		testName           string
//		email              string
//		expectedStatusCode int
//	}{
//		{testName: "existing email", email: user.Email, expectedStatusCode: http.StatusOK},
//		{testName: "invalid email address", email: "invalidemail", expectedStatusCode: http.StatusBadRequest},
//		{testName: "not existing email", email: "notinserted@mail.com", expectedStatusCode: http.StatusNotFound},
//	}
//
//	for _, tt := range sendRecoverTests {
//		t.Log(tt.testName)
//
//		uEntry := struct {
//			Email string `json:"email"`
//		}{
//			Email: tt.email,
//		}
//
//		body, _ := json.Marshal(uEntry)
//
//		rr := httptest.NewRecorder()
//		req := httptest.NewRequest(http.MethodPost, "/api/auth/recover", bytes.NewReader(body))
//		mainRouter.ServeHTTP(rr, req)
//
//		if rr.Code == http.StatusOK {
//			bodyRes := struct {
//				Link string `json:"link"`
//			}{}
//			if err := json.NewDecoder(rr.Body).Decode(&bodyRes); err != nil {
//				t.Error("could not decode the body response:", err.Error())
//			}
//			if bodyRes.Link == "" {
//				t.Error("link should not be empty if the response is OK")
//			}
//		}
//		if rr.Code != tt.expectedStatusCode {
//			t.Errorf("expected statusCode is %d, but got %d", tt.expectedStatusCode, rr.Code)
//		}
//	}
//}

func TestAuthController_ChangePass(t *testing.T) {
	password := "Password123"
	user := models.User{
		Name:      "Change Pass ctrl",
		Bio:       "This is the change pass controller",
		Email:     "changepass@controller.com",
		Password:  password,
		Phone:     8091234567,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, err := utils.GenerateHash(password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPass
	insertedID, err := authTestRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}
	token, err := json_web_token.Generate(insertedID, "")
	if err != nil {
		t.Error(err.Error())
	}

	var changePassTests = []struct {
		testName           string
		token              string
		newPassword        string
		expectedStatusCode int
	}{
		{testName: "valid token and password", token: token, newPassword: "ValidPass123456", expectedStatusCode: http.StatusOK},
		{testName: "invalid password format", token: token, newPassword: "asdasfaf", expectedStatusCode: http.StatusBadRequest},
		{testName: "invalid token", token: "invalid", newPassword: "ValidPass123456", expectedStatusCode: http.StatusUnauthorized},
	}

	for _, tt := range changePassTests {
		t.Log(tt.testName)
		uEntry := struct {
			NewPassword string `json:"newPassword"`
		}{
			NewPassword: tt.newPassword,
		}

		body, err := json.Marshal(uEntry)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewReader(body))
		req.Header.Add("Authorization", tt.token)
		mainRouter.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected status code is %d, but got %d", tt.expectedStatusCode, rr.Code)

			var resBody interface{}

			_ = json.NewDecoder(rr.Body).Decode(&resBody)
			t.Error("BODY RESPONSE", resBody)
		}
	}
}

//go:build integration
// +build integration

package services

import (
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
	"testing"
	"time"
)

func TestInitAuth(t *testing.T) {
	var i interface{}
	i = InitAuth(app, authRepo)
	if _, ok := i.(*Auth); !ok {
		t.Error("wrong type returned")
	}
}

func TestGetAuth(t *testing.T) {
	var i interface{}
	i = GetAuth()
	if _, ok := i.(*Auth); !ok {
		t.Error("wrong type returned")
	}
}

func TestAuth_Register(t *testing.T) {
	var registerTests = []struct {
		testName       string
		uEntry         models.User
		expectedErrors int
	}{
		{testName: "valid user", uEntry: models.User{
			Name:     "John Doe",
			Bio:      "I am john",
			Email:    "john@mail.com",
			Password: "HitsNow1224",
			Phone:    80976543212,
			PhotoURL: "",
		}, expectedErrors: 0},
		{testName: "duplicate email user", uEntry: models.User{
			Name:     "John Doe",
			Bio:      "I am john",
			Email:    "john@mail.com",
			Password: "HitsNow1224",
			Phone:    80976543212,
			PhotoURL: "",
		}, expectedErrors: 1},
	}

	for _, tt := range registerTests {
		t.Log(tt.testName)
		_, errColl := authTestService.Register(tt.uEntry)
		if len(errColl) != tt.expectedErrors {
			t.Errorf("expected %d but got %d", tt.expectedErrors, len(errColl))
		}
	}
}

func TestAuth_Login(t *testing.T) {
	password := "validPass1234"
	uEntry := models.User{
		Name:      "Login user",
		Bio:       "This is the login service user",
		Email:     "loginservice@mail.com",
		Password:  password,
		Phone:     1234567890,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashPass, err := utils.GenerateHash(uEntry.Password)
	if err != nil {
		t.Error(err.Error())
	}
	uEntry.Password = hashPass
	_, errColl := authTestService.Register(uEntry)
	if len(errColl) != 0 {
		t.Error(errColl)
	}

	var loginTests = []struct {
		testName, email, password string
		expectingErr              bool
	}{
		{testName: "valid user", email: uEntry.Email, password: password, expectingErr: false},
		{testName: "invalid email", email: "ramdon@dmm.com", password: password, expectingErr: true},
		{testName: "invalid password", email: uEntry.Email, password: "wrongpass", expectingErr: true},
	}

	for _, tt := range loginTests {
		t.Log(tt.testName)
		_, errColl := authTestService.Login(tt.email, tt.password)
		if len(errColl) != 0 && !tt.expectingErr {
			t.Error("was NOT expecting errors but got some")
		} else if len(errColl) == 0 && tt.expectingErr {
			t.Error("was expecting errors but did NOT get them")
		}
	}
}

func TestAuth_SendRecoverPassLink(t *testing.T) {
	password := "validPass1234"
	uEntry := models.User{
		Name:      "Send recover user",
		Bio:       "This is the send recover service",
		Email:     "recoverservice@mail.com",
		Password:  password,
		Phone:     1234567890,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashPass, err := utils.GenerateHash(uEntry.Password)
	if err != nil {
		t.Error(err.Error())
	}
	uEntry.Password = hashPass
	_, err = authRepo.Register(uEntry)
	if err != nil {
		t.Error(err.Error())
	}

	var sendRecoverTests = []struct {
		testName     string
		email        string
		expectingErr bool
	}{
		{testName: "valid email", email: uEntry.Email, expectingErr: false},
		{testName: "invalid email", email: "notvalid@mail.com", expectingErr: true},
	}

	for _, tt := range sendRecoverTests {
		t.Log(tt.testName)
		token, errSlice := authTestService.SendRecoverPassLink(tt.email)
		hasErrors := len(errSlice) != 0
		if hasErrors != tt.expectingErr {
			t.Errorf("expectingErr was %v, but got %v", tt.expectingErr, hasErrors)
		}

		if token == "" && !hasErrors {
			t.Error("token must not return empty if there are NOT errors")
		}
	}
}

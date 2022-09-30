package repositories_impl

import (
	"errors"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
	"testing"
	"time"
)

func TestInitAuthRepo(t *testing.T) {
	var i interface{}
	i = InitAuthRepo(app, dbConn)

	if _, ok := i.(*AuthRepoImpl); !ok {
		t.Error("wrong type")
	}
}

func TestGetAuthRepo(t *testing.T) {
	var i interface{}
	i = GetAuthRepo()
	if _, ok := i.(*AuthRepoImpl); !ok {
		t.Error("wrong type")
	}
}

func TestAuthRepoImpl_Register(t *testing.T) {
	var registerTests = []struct {
		testName      string
		uEntry        models.User
		expectedError bool
	}{
		{testName: "valid user", uEntry: models.User{
			Name:      "Auth test repo-1",
			Bio:       "Some random user",
			Email:     "valid-auth@mail.com",
			Password:  "SomePassword",
			Phone:     1234567898,
			PhotoURL:  "",
			IsActive:  true,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}, expectedError: false},
		{testName: "email in user", uEntry: models.User{
			Name:      "Auth test repo-2",
			Bio:       "Some random user 2",
			Email:     "valid-auth@mail.com",
			Password:  "SomePassword",
			Phone:     1234567898,
			PhotoURL:  "",
			IsActive:  true,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}, expectedError: true},
	}

	for _, tt := range registerTests {
		t.Log(tt.testName)
		hashedPass, err := utils.GenerateHash(tt.uEntry.Password)
		if err != nil {
			t.Error(err.Error())
		}
		tt.uEntry.Password = hashedPass
		_, err = authTestRepo.Register(tt.uEntry)

		if err != nil && !tt.expectedError {
			t.Error(err.Error())
		} else if err == nil && tt.expectedError {
			t.Error("expected error but got nothing")
		}
	}
}

func TestAuthRepoImpl_Login(t *testing.T) {
	password := "LoginPass1234"
	validUSer := models.User{
		Name:      "login user",
		Bio:       "this is a login user",
		Email:     "loginrepo@mail.com",
		Password:  password,
		Phone:     1234567890,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	hashPass, err := utils.GenerateHash(validUSer.Password)
	if err != nil {
		t.Error(err.Error())
	}
	validUSer.Password = hashPass

	_, err = authTestRepo.Register(validUSer)
	if err != nil {
		t.Error(err.Error())
	}

	var loginTests = []struct {
		testName        string
		email, password string
		expectedErr     error
	}{
		{testName: "valid credentials", email: validUSer.Email, password: password, expectedErr: nil},
		{testName: "invalid email", email: "some@mail", password: password, expectedErr: errors.New("wrong email or password")},
		{testName: "invalid password", email: validUSer.Email, password: "wrongpass", expectedErr: errors.New("wrong email or password")},
	}

	for _, tt := range loginTests {
		t.Log(tt.testName)
		_, err = authTestRepo.Login(tt.email, tt.password)
		if tt.expectedErr == nil && err != nil {
			t.Error(err.Error())
		} else if tt.expectedErr != nil && err == nil {
			t.Error("expected error is NOT 'nil' but did not get an error")
		}
	}
}

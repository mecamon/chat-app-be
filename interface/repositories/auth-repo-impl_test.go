//go:build integration
// +build integration

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
		{testName: "email in use", uEntry: models.User{
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

func TestAuthRepoImpl_FindByEmail(t *testing.T) {
	password := "validPassword12344"
	user := models.User{
		Name:      "user to FindByEmail",
		Bio:       "This is an user to find by email",
		Email:     "findby@mail.com",
		Password:  password,
		Phone:     8098907654,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashPassword

	_, err = authRepoImpl.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	var findByEmailTests = []struct {
		testName, email string
		userFound       bool
		expectedErr     error
	}{
		{testName: "user found", email: user.Email, userFound: true, expectedErr: nil},
		{testName: "user not found", email: "notexisting@mail.com", userFound: false, expectedErr: errors.New("has error")},
	}

	for _, tt := range findByEmailTests {
		t.Log(tt.testName)
		u, err := authTestRepo.FindByEmail(tt.email)

		if tt.expectedErr == nil {
			if err != nil {
				t.Error("error was not expected, but got error:", err.Error())
			}
			if u.Email != tt.email {
				t.Error("email inserted and email found are NOT the same")
			}
		} else if tt.expectedErr != nil {
			if err == nil {
				t.Error("error was expected but did NOT get any")
			}
		}
	}
}

func TestAuthRepoImpl_FindByID(t *testing.T) {
	password := "validPassword12344"
	user := models.User{
		Name:      "user to FindByID",
		Bio:       "This is an user to find by ID",
		Email:     "findby@id.com",
		Password:  password,
		Phone:     8098907654,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashPassword

	insertedID, err := authRepoImpl.Register(user)
	if err != nil {
		t.Error(err.Error())
	}
	notExistingID := "3eb3d668b31de5d588f42a6d"

	var findByIDTests = []struct {
		testName, uid string
		userFound     bool
		expectedErr   error
	}{
		{testName: "user found", uid: insertedID, userFound: true, expectedErr: nil},
		{testName: "user not found", uid: notExistingID, userFound: false, expectedErr: errors.New("has error")},
	}

	for _, tt := range findByIDTests {
		t.Log("TEST NAME:", tt.testName)
		u, err := authTestRepo.FindByID(tt.uid)

		if tt.expectedErr == nil {
			if err != nil {
				t.Error("error was not expected, but got error:", err.Error())
			}
			if u.ID.Hex() != tt.uid {
				t.Error("insertedID and user ID are NOT the same")
			}
		} else if tt.expectedErr != nil {
			t.Log("USER:", u)
			if err == nil {
				t.Error("error was expected but did NOT get any")
			}
		}
	}
}

func TestAuthRepoImpl_ChangePassword(t *testing.T) {
	password := "validPassword12344"
	user := models.User{
		Name:      "user to change",
		Bio:       "THis is an user to change the password",
		Email:     "change@mail.com",
		Password:  password,
		Phone:     8098907654,
		IsActive:  true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hasPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hasPassword
	insertedID, err := authRepoImpl.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	var changePassTests = []struct {
		testName, newPass, userID string
		expectedErr               error
	}{
		{testName: "successful changed password", newPass: "AnotherPass123", userID: insertedID, expectedErr: nil},
		{testName: "not existing ID", newPass: "AnotherPass123", userID: "notExistingId", expectedErr: errors.New("has error")},
	}

	for _, tt := range changePassTests {
		t.Log(tt.testName)
		err = authTestRepo.ChangePassword(tt.userID, tt.newPass)

		if tt.expectedErr == nil {
			if err != nil {
				t.Error("got an error,but was NOT expecting any. error:", err.Error())
			}
		}
		if tt.expectedErr != nil {
			if err == nil {
				t.Error("it was expecting an error but did NOT get it")
			}
		}
	}
}

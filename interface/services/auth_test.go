//go:build integration
// +build integration

package services

import (
	"github.com/mecamon/chat-app-be/models"
	"testing"
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

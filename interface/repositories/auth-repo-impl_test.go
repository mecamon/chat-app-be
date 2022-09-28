//go:build integration
// +build integration

package repositories_impl

import (
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/utils"
	"testing"
	"time"
)

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

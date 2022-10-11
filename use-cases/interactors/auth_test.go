//go:build !integration
// +build !integration

package interactors

import (
	"testing"

	"github.com/mecamon/chat-app-be/models"
)

func TestEvalRegistryEntry(t *testing.T) {
	var evalRegistryTests = []struct {
		testName       string
		userEntry      models.User
		expectedResult bool
		errorsExpected int
	}{
		{testName: "no name and no password", userEntry: models.User{
			Name:     "",
			Password: "",
			Email:    "valid@mail.com",
			Bio:      "This is a bio",
			Phone:    123456789,
		}, expectedResult: false, errorsExpected: 2},
		{testName: "wrong email format", userEntry: models.User{
			Name:     "John Doe",
			Password: "Password007",
			Email:    "not-validmailcom",
			Bio:      "This is a bio",
			Phone:    123456789,
		}, expectedResult: false, errorsExpected: 1},
		{testName: "wrong password format", userEntry: models.User{
			Name:     "Var Char",
			Password: "hbngddada",
			Email:    "valid@mail.com",
			Bio:      "This is a bio",
			Phone:    123456789,
		}, expectedResult: false, errorsExpected: 1},
		{testName: "phone not long enough", userEntry: models.User{
			Name:     "Random name",
			Password: "PasswordValid009",
			Email:    "valid@mail.com",
			Bio:      "This is a bio",
			Phone:    1234567,
		}, expectedResult: false, errorsExpected: 1},
		{testName: "valid user entry", userEntry: models.User{
			Name:     "Random name",
			Password: "PasswordValid009",
			Email:    "valid@mail.com",
			Bio:      "This is a bio",
			Phone:    123456789,
		}, expectedResult: true, errorsExpected: 0},
	}

	for _, tt := range evalRegistryTests {
		t.Log(tt.testName)
		result, errors := EvalRegistryEntry(tt.userEntry)
		if result != tt.expectedResult {
			t.Errorf("expected result is: %v, but got %v", tt.expectedResult, result)
		}
		if len(errors) != tt.errorsExpected {
			t.Errorf("expected errors are: %d, but got %d", tt.errorsExpected, len(errors))
		}
	}
}

func TestCompleteRegEntry(t *testing.T) {
	uEntry := models.User{
		Name:     "Carlos",
		Bio:      "I am a web developer",
		Email:    "valid@mail.com",
		Password: "PAssword008",
		Phone:    98765432134,
	}

	completedUEntry := CompleteRegEntry(uEntry)
	if completedUEntry.CreatedAt == 0 || completedUEntry.UpdatedAt == 0 || !completedUEntry.IsActive {
		t.Error("expected CreatedAt and UpdatedAt to be defined but they aren't")
	}
}

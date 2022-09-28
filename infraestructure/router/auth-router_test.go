package router

import (
	"bytes"
	"encoding/json"
	"github.com/mecamon/chat-app-be/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthController_Register(t *testing.T) {
	var registerTests = []struct {
		testName     string
		uEntry       models.User
		expectedCode int
	}{
		{testName: "valid user credentials", uEntry: models.User{
			Name:     "Juan Lopez",
			Bio:      "Yo soy Juan",
			Email:    "juan@mail.com",
			Password: "JuanMan123456",
			Phone:    8097654312,
		}, expectedCode: http.StatusCreated},
		{testName: "email in use", uEntry: models.User{
			Name:     "Juan Lopez",
			Bio:      "Yo soy Juan",
			Email:    "juan@mail.com",
			Password: "JuanMan123456",
			Phone:    8097654312,
		}, expectedCode: http.StatusBadRequest},
		{testName: "no name", uEntry: models.User{
			Name:     "",
			Bio:      "Yo soy Juan",
			Email:    "juan@mail.com",
			Password: "JuanMan123456",
			Phone:    8097654312,
		}, expectedCode: http.StatusBadRequest},
		{testName: "wrong formatted email and password", uEntry: models.User{
			Name:     "",
			Bio:      "Yo soy Juan",
			Email:    "juanmail.com",
			Password: "12345678676556",
			Phone:    8097654312,
		}, expectedCode: http.StatusBadRequest},
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
	}
}
//go:build integration
// +build integration

package services

import (
	"testing"
)

func TestInitMailService(t *testing.T) {
	var i interface{}
	i = InitMailService(app)

	if _, ok := i.(*Mail); !ok {
		t.Error("wrong type")
	}
}

func TestGetMail(t *testing.T) {
	var i interface{}
	i = GetMail()
	if _, ok := i.(*Mail); !ok {
		t.Error("wrong type")
	}
}

//
//func TestEmail_SendRecoverPassLink(t *testing.T) {
//	info := models.EmailInfo{
//		Address: app.EmailAcc,
//		Subject: "Test mail",
//		Body:    fmt.Sprint("<h4>You are testing the password recovery</h4>"),
//	}
//	if err := mailTestService.Recovery(info); err != nil {
//		t.Error("error sending mail:", err.Error())
//	}
//}

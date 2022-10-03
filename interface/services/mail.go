package services

import (
	"github.com/mecamon/chat-app-be/config"
)

type Mail struct {
	app *config.App
}

var mail *Mail

func InitMailService(app *config.App) *Mail {
	mail := &Mail{app: app}
	return mail
}

func GetMail() *Mail {
	return mail
}

//func (e *Mail) Recovery(info models.EmailInfo) error {
//	m := gomail.NewMessage()
//	m.SetHeader("From", e.app.EmailAcc)
//	m.SetHeader("To", info.Address)
//	m.SetHeader("Subject", info.Subject)
//	m.SetBody("text/html", info.Body)
//
//	d := gomail.NewDialer(
//		e.app.EmailHost,
//		e.app.EmailPort,
//		e.app.EmailAcc,
//		e.app.EmailAccPass)
//
//	// This is only needed when SSL/TLS certificate is not valid on server.
//	// In production this should be set to false.
//	if !e.app.IsProd {
//		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
//	}
//
//	if err := d.DialAndSend(m); err != nil {
//		return err
//	}
//	return nil
//}

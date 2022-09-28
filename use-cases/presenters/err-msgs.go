package presenters

import (
	appi18n "github.com/mecamon/chat-app-be/i18n"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
)

func ErrMessages(loc appi18n.AppLocales, ce []*cErrors.Custom) []string {
	var errMessages []string

	for _, e := range ce {
		message := loc.GetMsg(e.MessageID, e.TemplateData)
		e.SetLocalesErrMsg(message)
		errMessages = append(errMessages, message)
	}

	return errMessages
}

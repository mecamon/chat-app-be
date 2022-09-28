//go:build !integration
// +build !integration

package presenters

import (
	appi18n "github.com/mecamon/chat-app-be/i18n"
	cErrors "github.com/mecamon/chat-app-be/use-cases/c-errors"
	"testing"
)

func TestErrMessages(t *testing.T) {
	if err := appi18n.InitLocales(); err != nil {
		t.Error(err.Error())
	}
	multi := appi18n.GetMultiLocales()
	loc := multi.GetSpeLocales("en-EN")

	errors := []*cErrors.Custom{
		{
			Property:     "name",
			MessageID:    "NameTooShort",
			TemplateData: map[string]interface{}{"Count": 2},
		},
		{
			Property:     "email",
			MessageID:    "EmailAddressTaken",
			TemplateData: nil,
		},
	}

	result := ErrMessages(loc, errors)
	if len(result) != 2 {
		t.Error("expected 2 but got less or more")
	}

}

//go:build integration
// +build integration

package data

import (
	"context"
	"github.com/mecamon/chat-app-be/config"
	"os"
	"testing"
)

var app *config.App
var dbConn *DB

func TestMain(m *testing.M) {
	config.SetConfig()
	app = config.GetConfig()

	code := m.Run()
	shutdown(dbConn)
	os.Exit(code)
}

func shutdown(dbConn *DB) {
	if err := dbConn.Client.Disconnect(context.TODO()); err != nil {
		panic(err.Error())
	}
}

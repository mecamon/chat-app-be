//go:build integration
// +build integration

package data

import (
	"fmt"
	"testing"
)

func TestCreateDBClient(t *testing.T) {
	var err error

	dbConnUri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?maxPoolSize=20&w=majority",
		app.DBUser,
		app.DBUserPassword,
		app.DBHost,
		app.DBPort,
		app.DBName)
	dbConn, err = CreateDBClient(dbConnUri)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestPingDB(t *testing.T) {
	if err := PingDB(dbConn.Client); err != nil {
		t.Error(err.Error())
	}
}

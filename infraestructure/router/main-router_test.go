//go:build integration
// +build integration

package router

import "testing"

func TestGetMain(t *testing.T) {
	var i interface{}
	SetRouter()
	i, err := GetMain()
	if err != nil {
		t.Error(err.Error())
	}

	if _, ok := i.(*Main); !ok {
		t.Error("wrong type returned")
	}
}

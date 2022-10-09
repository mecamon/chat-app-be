//go:build !integration
// +build !integration

package utils

import (
	"fmt"
	"testing"
)

func TestGetRouterParam(t *testing.T) {
	protocol := "http://"
	domain := "example.com/"
	path := "items/"
	routeParam := "123"
	uri := fmt.Sprintf("%s%s%s%s", protocol, domain, path, routeParam)

	p := GetRouteParam(uri)
	if routeParam != p {
		t.Errorf("expected route param is %s, but got %s instead", routeParam, p)
	}
}
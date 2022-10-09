package utils

import "strings"

func GetRouteParam(uri string) string {
	parts := strings.Split(uri, "/")
	lastI := len(parts) - 1
	rParam := parts[lastI]
	return rParam
}

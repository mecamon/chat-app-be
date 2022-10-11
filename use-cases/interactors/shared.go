package interactors

import (
	"github.com/mecamon/chat-app-be/models"
	"regexp"
)

func ValidFile(fileInfo models.FileInfo, maxSize int64, contentType ...string) bool {
	if fileInfo.Size > maxSize {
		return false
	}
	var pattern string

	for i, c := range contentType {
		if i == 0 {
			pattern += c
		} else {
			pattern += "|"
			pattern += c
		}
	}

	reg, _ := regexp.Compile(pattern)
	return reg.MatchString(fileInfo.ContentType)
}

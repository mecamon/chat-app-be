//go:build !integration
// +build !integration

package interactors

import (
	"github.com/mecamon/chat-app-be/models"
	"testing"
)

func TestValidFile(t *testing.T) {
	const maxFileSize = 5242880

	var validFileTests = []struct {
		testName       string
		fileInfo       models.FileInfo
		expectedResult bool
		contentTypes   []string
	}{
		{testName: "invalid file size", fileInfo: models.FileInfo{
			Size:        12344444,
			ContentType: "image/jpg",
		}, expectedResult: false, contentTypes: []string{"image/jpg", "image/png"}},
		{testName: "invalid file content type", fileInfo: models.FileInfo{
			Size:        maxFileSize - 500,
			ContentType: "document/pdf",
		}, expectedResult: false, contentTypes: []string{"image/jpg", "image/png"}},
		{testName: "valid file", fileInfo: models.FileInfo{
			Size:        maxFileSize - 1000,
			ContentType: "image/jpg",
		}, expectedResult: true, contentTypes: []string{"image/jpg", "image/png"}},
	}

	for _, tt := range validFileTests {
		t.Log("TEST NAME:", tt.testName)
		matched := ValidFile(tt.fileInfo, maxFileSize, tt.contentTypes...)
		if matched != tt.expectedResult {
			t.Errorf("expected result is %v, but got %v instead", tt.expectedResult, matched)
		}
	}
}

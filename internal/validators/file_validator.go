package validators

import (
	"fmt"
	"mime/multipart"
	"strings"
)

type FileValidator interface {
	ValidateFileSize(file multipart.FileHeader) error
	ValidateFileType(file multipart.FileHeader) error
}

func ValidateFileSize(file multipart.FileHeader) error {
	if file.Size > 100000 {
		return fmt.Errorf("File size exceeds the limit, maximum size is 100KB")
	}

	return nil
}

func ValidateFileType(file multipart.FileHeader) error {
	allowedTypes := []string{"image/jpeg", "image/png", "image/jpg"}
	contentType := file.Header.Get("Content-Type")

	isValidType := false

	if contentType == "application/octet-stream" {
		if strings.Contains(file.Filename, ".jpg") ||
			strings.Contains(file.Filename, ".jpeg") ||
			strings.Contains(file.Filename, ".png") {
			isValidType = true
		}
	} else {
		for _, t := range allowedTypes {
			if contentType == t {
				isValidType = true
				break
			}
		}
	}

	if !isValidType {
		return fmt.Errorf("Mime type %s is not allowed", contentType)
	}

	return nil
}

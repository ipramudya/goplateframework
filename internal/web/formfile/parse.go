package formfile

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
)

func Parse(file *multipart.FileHeader, givenType string) ([]byte, error) {
	contentType := file.Header.Get("Content-Type")

	if contentType == "" || strings.Contains(contentType, givenType) {
		return nil, fmt.Errorf("image: content type %s is not supported", contentType)
	}

	source, err := file.Open()
	if err != nil {
		return nil, errors.New("parse: file cannot be opened")
	}
	defer source.Close()

	data := make([]byte, file.Size)
	if _, err := source.Read(data); err != nil {
		return nil, errors.New("parse: file unreadable")
	}

	return data, nil
}

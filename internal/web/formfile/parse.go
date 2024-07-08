package formfile

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

func Parse(file *multipart.FileHeader, desiredContentType string) ([]byte, error) {
	source, err := file.Open()
	if err != nil {
		return nil, errors.New("parse: file cannot be opened")
	}
	defer source.Close()

	data := make([]byte, file.Size)
	if _, err := source.Read(data); err != nil {
		return nil, errors.New("parse: file unreadable")
	}

	contentType := http.DetectContentType(data)

	if strings.HasSuffix(desiredContentType, "/*") {
		if !strings.HasPrefix(contentType, strings.TrimSuffix(desiredContentType, "*")) {
			return nil, fmt.Errorf("parse: content type %s does not match %s", contentType, desiredContentType)
		}
	} else if contentType != desiredContentType {
		return nil, fmt.Errorf("parse: content type %s does not match %s", contentType, desiredContentType)
	}

	return data, nil
}

package formfile

import "mime/multipart"

func Parse(file *multipart.FileHeader) ([]byte, error) {
	source, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()

	data := make([]byte, file.Size)
	if _, err := source.Read(data); err != nil {
		return nil, err
	}

	return data, nil
}

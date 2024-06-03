package validate

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var Phone = validation.By(func(value interface{}) error {
	s, _ := value.(string)
	if !strings.HasPrefix(s, "+62") {
		return errors.New("must start with +62")
	}
	return nil
})

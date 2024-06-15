package validate

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	Phone = validation.By(func(value interface{}) error {
		s, _ := value.(string)
		if !strings.HasPrefix(s, "+62") {
			return errors.New("must start with +62")
		}
		return nil
	})

	Timestamp = validation.By(func(value interface{}) error {
		s, _ := value.(string)
		if _, err := time.Parse(time.RFC3339, s); err != nil {
			return errors.New("invalid timestamp")
		}
		return nil
	})
)

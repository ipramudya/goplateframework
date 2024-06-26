package menu

import validation "github.com/go-ozzo/ozzo-validation/v4"

type MenuDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	ImageURL    string  `json:"image_url"`
}

type NewMenuDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	// ImageURL    string `json:"image_url"`
	OutletID string `json:"outlet_id"`
}

func (m NewMenuDTO) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.Description, validation.Required, validation.Length(1, 255)),
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.OutletID, validation.Required),
	)
}

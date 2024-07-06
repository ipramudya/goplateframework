package addressweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/sdk/validate"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

// required usecase methods which this controller needs to operate the business logic
type iUsecase interface {
	Update(ctx context.Context, na *address.NewAddressDTO, id uuid.UUID) (*address.AddressDTO, error)
}

type controller struct {
	addressUC iUsecase
	log       *logger.Log
}

func newController(addressUC iUsecase, log *logger.Log) *controller {
	return &controller{addressUC, log}
}

func (con *controller) update(c echo.Context) error {
	na := new(address.NewAddressDTO)

	if err := c.Bind(na); err != nil {
		return errshttp.New(errshttp.InvalidArgument, "Given JSON is invalid")
	}

	if err := na.Validate(); err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given JSON is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Address id is invalid, should be valid UUID")
		e.AddDetail("id: invalid")
		return e
	}

	a, err := con.addressUC.Update(c.Request().Context(), na, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}

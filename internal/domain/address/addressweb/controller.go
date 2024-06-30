package addressweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/sdk/errs"
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
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := na.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid id: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	a, err := con.addressUC.Update(c.Request().Context(), na, id)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, a)
}

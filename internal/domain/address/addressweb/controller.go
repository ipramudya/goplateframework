package addressweb

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/domain/address/addressuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type controller struct {
	addressUC *addressuc.Usecase
	log       *logger.Log
}

func newController(addressUC *addressuc.Usecase, log *logger.Log) *controller {
	return &controller{addressUC, log}
}

func (con *controller) update(c echo.Context) error {
	dto := new(address.NewAddressDTO)

	if err := c.Bind(dto); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := dto.Validate(); err != nil {
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

	a, err := con.addressUC.Update(c.Request().Context(), dto, id.String())

	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, a)
}

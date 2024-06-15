package outletdelivery

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/outlet"
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type controller struct {
	outletUC *outletuc.Usecase
	log      *logger.Log
}

func newController(outletUC *outletuc.Usecase, log *logger.Log) *controller {
	return &controller{outletUC, log}
}

func (con *controller) create(c echo.Context) error {
	dto := new(outlet.NewOutletDTO)

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

	outlet, err := con.outletUC.Create(c.Request().Context(), dto)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, outlet)
}

func (con *controller) getOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		e := errs.New(errs.InvalidArgument, errors.New("invalid id"))
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	outlet, err := con.outletUC.GetOne(c.Request().Context(), id.String())
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, outlet)
}

func (con *controller) update(c echo.Context) error {
	dto := new(outlet.NewOutletDTO)

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

	outlet, err := con.outletUC.Update(c.Request().Context(), dto, id.String())
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, outlet)
}

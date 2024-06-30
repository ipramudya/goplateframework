package outletxweb

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/addressx"
	"github.com/goplateframework/internal/domain/outletx"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type iAddressUsecase interface {
	Create(ctx context.Context, na *addressx.NewAddressDTO) (*addressx.AddressDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*addressx.AddressDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type iOutletUsecase interface {
	Create(ctx context.Context, no *outletx.NewOutletDTO) (*outletx.OutletDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*outletx.OutletDTO, error)
	Update(ctx context.Context, no *outletx.NewOutletDTO, id uuid.UUID) (*outletx.OutletDTO, error)
}

type controller struct {
	addressUC iAddressUsecase
	outletUC  iOutletUsecase
	log       *logger.Log
}

func newController(addressUC iAddressUsecase, outletUC iOutletUsecase, log *logger.Log) *controller {
	return &controller{addressUC, outletUC, log}
}

func (con *controller) create(c echo.Context) error {
	no := new(outletx.NewOutletDTO)

	if err := c.Bind(no); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := no.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	a, err := con.addressUC.Create(c.Request().Context(), &addressx.NewAddressDTO{
		Street:     no.Address.Street,
		City:       no.Address.City,
		Province:   no.Address.Province,
		PostalCode: no.Address.PostalCode,
	})

	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	no.Address = a
	o, err := con.outletUC.Create(c.Request().Context(), no)
	if err != nil {
		if err := con.addressUC.Delete(c.Request().Context(), a.ID); err != nil {
			return c.JSON(err.(*errs.Error).HTTPStatus(), err)
		}

		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, o)
}

func (con *controller) getOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		e := errs.New(errs.InvalidArgument, errors.New("invalid id"))
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	o, err := con.outletUC.GetOne(c.Request().Context(), id)

	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, o)
}

func (con *controller) update(c echo.Context) error {
	no := new(outletx.NewOutletDTO)

	if err := c.Bind(no); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := no.Validate(); err != nil {
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

	o, err := con.outletUC.Update(c.Request().Context(), no, id)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, o)
}

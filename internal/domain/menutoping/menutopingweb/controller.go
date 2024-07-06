package menutopingweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/menutoping"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/sdk/validate"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type iUsecase interface {
	Create(ctx context.Context, nmt *menutoping.NewMenuTopingsDTO) (*menutoping.MenuTopingsDTO, error)
	GetAll(ctx context.Context) ([]*menutoping.MenuTopingsDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*menutoping.MenuTopingsDTO, error)
	Update(ctx context.Context, nmt *menutoping.NewMenuTopingsDTO, id uuid.UUID) (*menutoping.MenuTopingsDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type controller struct {
	menuTopingUC iUsecase
	log          *logger.Log
}

func newController(menuTopingUC iUsecase, log *logger.Log) *controller {
	return &controller{
		menuTopingUC: menuTopingUC,
		log:          log,
	}
}

func (con *controller) create(c echo.Context) error {
	nmt := new(menutoping.NewMenuTopingsDTO)

	if err := c.Bind(nmt); err != nil {
		return errshttp.New(errshttp.InvalidArgument, "Given JSON is invalid")
	}

	if err := nmt.Validate(); err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given JSON is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	m, err := con.menuTopingUC.Create(c.Request().Context(), nmt)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func (con *controller) getAll(c echo.Context) error {
	m, err := con.menuTopingUC.GetAll(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) getOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Menu topings id is invalid, should be valid UUID")
		e.AddDetail("id: invalid")
		return e
	}

	m, err := con.menuTopingUC.GetOne(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) update(c echo.Context) error {
	nmt := new(menutoping.NewMenuTopingsDTO)

	if err := c.Bind(nmt); err != nil {
		return errshttp.New(errshttp.InvalidArgument, "Given JSON is invalid")
	}

	if err := nmt.Validate(); err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given JSON is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Menu topings id is invalid, should be valid UUID")
		e.AddDetail("id: invalid")
		return e
	}

	m, err := con.menuTopingUC.Update(c.Request().Context(), nmt, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Menu topings id is invalid, should be valid UUID")
		e.AddDetail("id: invalid")
		return e
	}

	err = con.menuTopingUC.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

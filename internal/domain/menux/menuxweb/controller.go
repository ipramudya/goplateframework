package menuxweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/goplateframework/internal/domain/menux"
)

// required usecase methods which this controller needs to operate the business logic
type Usecase interface {
	Create(ctx context.Context, nm *menux.NewMenuDTO) (*menux.MenuDTO, error)
	GetAll(ctx context.Context, qp *QueryParams) (*[]menux.MenuDTO, error)
	Update(ctx context.Context, nm *menux.NewMenuDTO, id string) (*menux.MenuDTO, error)
}

type controller struct {
	menuUC Usecase
	log    *logger.Log
}

func newController(menuUC Usecase, log *logger.Log) *controller {
	return &controller{menuUC, log}
}

func (con *controller) create(c echo.Context) error {
	nm := new(menux.NewMenuDTO)

	if err := c.Bind(nm); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := nm.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	m, err := con.menuUC.Create(c.Request().Context(), nm)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, m)
}

func (con *controller) getAll(c echo.Context) error {
	qp, err := getQueryParams(c).Parse(c)

	if err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	m, err := con.menuUC.GetAll(c.Request().Context(), qp)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) update(c echo.Context) error {
	nm := new(menux.NewMenuDTO)

	if err := c.Bind(nm); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := nm.Validate(); err != nil {
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

	m, err := con.menuUC.Update(c.Request().Context(), nm, id.String())
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, m)
}

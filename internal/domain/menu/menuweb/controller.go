package menuweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/goplateframework/internal/domain/menu"
)

// required usecase methods which this controller needs to operate the business logic
type iUsecase interface {
	Create(ctx context.Context, nm *menu.NewMenuDTO) (*menu.MenuDTO, error)
	GetAll(ctx context.Context, qp *QueryParams) (*result.Result[menu.MenuDTO], error)
	Update(ctx context.Context, nm *menu.NewMenuDTO, id uuid.UUID) (*menu.MenuDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type controller struct {
	menuUC iUsecase
	log    *logger.Log
}

func newController(menuUC iUsecase, log *logger.Log) *controller {
	return &controller{menuUC, log}
}

func (con *controller) create(c echo.Context) error {
	nm := new(menu.NewMenuDTO)

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
	qp, err := getQueryParams(c).Parse()

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
	nm := new(menu.NewMenuDTO)

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

	m, err := con.menuUC.Update(c.Request().Context(), nm, id)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) delete(c echo.Context) error {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid id: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	err = con.menuUC.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.NoContent(http.StatusOK)
}

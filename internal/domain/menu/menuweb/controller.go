package menuweb

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/menu"
	"github.com/goplateframework/internal/domain/menu/menuuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type controller struct {
	menuUC *menuuc.Usecase
	log    *logger.Log
}

func newController(menuUC *menuuc.Usecase, log *logger.Log) *controller {
	return &controller{menuUC, log}
}

func (con *controller) create(c echo.Context) error {
	dto := new(menu.NewMenuDTO)

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

	menu, err := con.menuUC.Create(c.Request().Context(), dto)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, menu)
}

func (con *controller) update(c echo.Context) error {
	dto := new(menu.NewMenuDTO)

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

	menu, err := con.menuUC.Update(c.Request().Context(), dto, id.String())
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, menu)
}

func (con *controller) getAllByOutletID(c echo.Context) error {
	qp, err := menu.ParseQueryParams(c)
	if err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	menus, err := con.menuUC.GetAllByOutletID(c.Request().Context(), qp)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, menus)
}

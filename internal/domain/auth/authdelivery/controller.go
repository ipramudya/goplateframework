package authdelivery

import (
	"net/http"

	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/labstack/echo/v4"
)

type controller struct {
	authUC *authuc.Usecase
}

func newController(authUC *authuc.Usecase) *controller {
	return &controller{authUC}
}

func (con *controller) login(c echo.Context) error {
	dto := new(auth.LoginDTO)

	if err := c.Bind(dto); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := dto.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		return c.JSON(e.HTTPStatus(), e)
	}

	account, err := con.authUC.Login(c.Request().Context(), dto.Email, dto.Password)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, account)
}

func (con *controller) logout(c echo.Context) error {
	err := con.authUC.Logout(c.Request().Context())
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.NoContent(http.StatusNoContent)
}

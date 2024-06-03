package accountdelivery

import (
	"net/http"

	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/labstack/echo/v4"
)

type controller struct {
	accountUC *accountuc.Usecase
}

func newController(accountUC *accountuc.Usecase) *controller {
	return &controller{accountUC}
}

func (con *controller) register(c echo.Context) (err error) {
	dto := new(account.NewAccouuntDTO)

	if err := c.Bind(dto); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := dto.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		return c.JSON(e.HTTPStatus(), e)
	}

	account, err := con.accountUC.Register(c.Request().Context(), dto)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, account)
}

func (con *controller) changePassword(c echo.Context) error {
	dto := new(account.ChangePasswordDTO)

	if err := c.Bind(dto); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := dto.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		return c.JSON(e.HTTPStatus(), e)
	}

	err := con.accountUC.ChangePassword(c.Request().Context(), dto.OldPassword, dto.NewPassword)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.NoContent(http.StatusAccepted)
}

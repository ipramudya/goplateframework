package menuweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/sdk/validate"
	"github.com/goplateframework/internal/web/formfile"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/goplateframework/internal/domain/menu"
)

// required usecase methods which this controller needs to operate the business logic
type iUsecase interface {
	Create(ctx context.Context, nm *menu.NewMenuDTO, image *[]byte) (*menu.MenuDTO, error)
	GetAll(ctx context.Context, qp *QueryParams) (*result.Result[menu.MenuDTO], error)
	Update(ctx context.Context, nm *menu.NewMenuDTO, id uuid.UUID, image *[]byte) (*menu.MenuDTO, error)
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
		return errshttp.New(errshttp.InvalidArgument, "Given form-data is invalid")
	}

	if err := nm.Validate(); err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given form-data is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	file, err := c.FormFile("image")
	if err != nil {
		return errshttp.New(errshttp.InvalidArgument, "Given image form-data is invalid")
	}

	if file == nil {
		return errshttp.New(errshttp.InvalidArgument, "Image is required")
	}

	menuImage, err := formfile.Parse(file, "image/*")
	if err != nil {
		e := errshttp.New(errshttp.Internal, "Cannot parse given file")
		e.AddDetail(err.Error())
		return e
	}

	m, err := con.menuUC.Create(c.Request().Context(), nm, &menuImage)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func (con *controller) getAll(c echo.Context) error {
	qp, err := getQueryParams(c).Parse()

	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given query params are invalid")
		e.AddDetail(err.Error())
		return e
	}

	m, err := con.menuUC.GetAll(c.Request().Context(), qp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) update(c echo.Context) error {
	nm := new(menu.NewMenuDTO)

	if err := c.Bind(nm); err != nil {
		return errshttp.New(errshttp.InvalidArgument, "Given form-data is invalid")
	}

	if err := nm.Validate(); err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given form-data is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	if nm.ImageURL == "" {
		e := errshttp.New(errshttp.InvalidArgument, "Given form-data is out of validation rules")
		e.AddDetail("image_url: cannot be blank")
		return e
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Menu id is invalid, should be valid UUID")
		e.AddDetail("id: invalid")
		return e
	}

	file, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return errshttp.New(errshttp.InvalidArgument, "Given image form-data is invalid")
	}

	var menuImage *[]byte

	if file != nil {
		mi, err := formfile.Parse(file, "image/*")
		if err != nil {
			e := errshttp.New(errshttp.Internal, "Cannot parse given file")
			e.AddDetail(err.Error())
			return e
		}
		menuImage = &mi
	}

	m, err := con.menuUC.Update(c.Request().Context(), nm, id, menuImage)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (con *controller) delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Menu id is invalid, should be valid UUID")
		e.AddDetail("id: invalid")
		return e
	}

	err = con.menuUC.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

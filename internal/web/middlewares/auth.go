package middlewares

import (
	"errors"

	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/jsonwebtoken"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/labstack/echo/v4"
)

func (mid *Middleware) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		claims, err := jsonwebtoken.Validate(mid.conf, authHeader)

		if err != nil {
			if errors.Is(err, jsonwebtoken.ErrInvalidToken) {
				e := errs.Newf(errs.Unauthenticated, "unauthenticated: %v", err)
				mid.log.Debug(e.Debug())
				return c.JSON(e.HTTPStatus(), e)
			}

			e := errs.New(errs.Internal, err)
			mid.log.Debug(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		ctx := webcontext.SetClaims(c.Request().Context(), claims)
		ctx = webcontext.SetAccountPayload(ctx, &claims.Payload)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

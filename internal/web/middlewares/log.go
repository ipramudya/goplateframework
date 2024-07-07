package middlewares

import (
	"net/http"
	"time"

	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/labstack/echo/v4"
)

func (mid *Middleware) RequestLoggerMware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)

		req := c.Request()
		res := c.Response()

		since := time.Since(start).String()
		reqId := res.Header().Get(echo.HeaderXRequestID)

		mid.log.Infof("Req ID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s",
			reqId,
			req.Method,
			req.RequestURI,
			res.Status,
			res.Size,
			since,
		)

		return err
	}
}

func (mid *Middleware) ErrorLoggingMware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			switch e := err.(type) {
			case *errshttp.ErrorResponse:
				reqId := c.Response().Header().Get(echo.HeaderXRequestID)

				e.AddRequestID(reqId)
				mid.log.Error(e.LogForDebug())

				return c.JSON(e.HTTPStatus(), e)
			default:
				mid.log.Error(err.Error())
				return c.JSON(http.StatusInternalServerError, struct {
					Error string `json:"error"`
				}{
					Error: err.Error(),
				})
			}
		}

		return nil
	}
}

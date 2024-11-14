package isulog

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pirosiki197/isulog/internal"
)

func Echo() echo.MiddlewareFunc {
	return EchoWithConfig(DefaultConfig)
}

func EchoWithConfig(config Config) echo.MiddlewareFunc {
	recorder := internal.NewRecorder(config.Filename)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			if err := next(c); err != nil {
				return err
			}

			r := internal.Record{
				Path:         c.Path(),
				Method:       c.Request().Method,
				StatusCode:   c.Response().Status,
				ResponseTime: time.Since(start),
			}
			if err := recorder.Save(r); err != nil {
				c.Logger().Error(err)
			}

			return nil
		}
	}

}

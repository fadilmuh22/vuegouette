package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/fadilmuh22/restskuy/internal/handler"
)

func TransformErrorResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				return handler.SendResponse(c, he.Code, false, he.Message.(string), nil)
			}

			return handler.SendResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
		}

		return nil
	}
}

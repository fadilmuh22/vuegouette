package middleware

import (
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"

	"github.com/labstack/echo/v4"
)

func TransformErrorResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.BasicResponse{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		}

		return nil
	}
}

package handler

import (
	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	HandleRoutes(g *echo.Group)
}

func SendResponse(c echo.Context, status int, success bool, message string, data interface{}) error {
	return c.JSON(status, &model.BasicResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

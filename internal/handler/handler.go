package handler

import (
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func NewApiHandlers(e *echo.Echo, db *gorm.DB) {
	api := e.Group("api")

	NewStaticHandler().HandleRoutes(api)
	NewUserHandler(db).HandleRoutes(api)
	NewProductHandler(db).HandleRoutes(api)
	NewAuthHandler(db).HandleRoutes(api)
}

package handler

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler interface {
	HandleRoutes(g *echo.Group)
}

func NewApiHandlers(e *echo.Echo, db *gorm.DB) {
	api := e.Group("api")

	NewStaticHandler().HandleRoutes(api)
	NewUserHandler(db).HandleRoutes(api)
	NewAuthHandler(db).HandleRoutes(api)
	NewVideoHandler(db).HandleRoutes(api)
}

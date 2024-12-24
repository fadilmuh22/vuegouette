package handler

import (
	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler interface {
	HandleRoutes(g *echo.Group)
}

func NewApiHandlers(e *echo.Echo, db *gorm.DB, redisClient *db.RedisClient) {
	// INFO: WEB UI
	e.Static("", "static/")
	
	api := e.Group("api")

	NewStaticHandler().HandleRoutes(api)
	NewUserHandler(db).HandleRoutes(api)
	NewAuthHandler(db).HandleRoutes(api)
	NewVideoHandler(db, redisClient).HandleRoutes(api)
}

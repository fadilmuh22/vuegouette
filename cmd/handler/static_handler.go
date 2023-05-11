package handler

import "github.com/labstack/echo/v4"

type staticHandler struct{}

func (h staticHandler) HandleRoutes(g *echo.Group) {
	g.Static("/", "static/swaggerui")
	g.File("/swagger.yaml", "static/swagger.yaml")
}

func NewStaticHandler() Handler {
	return staticHandler{}
}

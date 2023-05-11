package handler

import "github.com/labstack/echo/v4"

type StaticHandler struct {
}

func (h StaticHandler) HandleRoutes(g *echo.Group) {
	g.Static("/", "static/swaggerui")
	g.File("/swagger.yaml", "static/swagger.yaml")
}

func NewStaticHandler() StaticHandler {
	return StaticHandler{}
}

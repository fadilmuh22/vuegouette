package main

import (
	"restskuy/cmd/service/user"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func setupRouter(e *echo.Echo) {

	api := e.Group("/api")

	user.HandleUser(api)

}

// go server using echo
func StartServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	setupRouter(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

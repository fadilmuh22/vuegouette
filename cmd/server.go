package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fadilmuh22/restskuy/cmd/service/product"
	"github.com/fadilmuh22/restskuy/cmd/service/user"
	"github.com/fadilmuh22/restskuy/cmd/utils"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func setupRouter(e *echo.Echo) {
	api := e.Group("api")

	e.Static("/", "static/swaggerui")
	e.File("/swagger.yaml", "static/swagger.yaml")

	user.HandleRoutes(api)
	product.HandleRoutes(api)
}

// echo middleware for transforming response to restful structure using struct called BasicResponse

// go server using echo
func StartServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// e.Use(middleware.TransformErrorResponse)

	cv := utils.CustomValidator{}
	err := cv.Init()
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Validator = &cv

	// Routes
	setupRouter(e)

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

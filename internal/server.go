package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/handler"
	"github.com/fadilmuh22/restskuy/internal/middleware"
	"github.com/fadilmuh22/restskuy/internal/util"
)

func runServer(e *echo.Echo) {
	go func() {
		var err error

		if viper.GetString("ENV") == "production" {
			err = e.StartAutoTLS(":1323")
		} else {
			err = e.Start(":1323")
		}

		if err != nil && err != http.ErrServerClosed {
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

// go server using echo
func StartServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	if viper.GetString("ENV") == "production" {
		e.Use(echomiddleware.Secure())
		e.Use(echomiddleware.HTTPSRedirect())
		e.Use(echomiddleware.HTTPSWWWRedirect())
	}

	e.Use(middleware.TransformErrorResponse)

	// Validator
	cv := util.CustomValidator{}
	err := cv.Init()
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Validator = &cv

	db := db.Connect()

	// Routes
	handler.NewApiHandlers(e, db)

	// Start server
	runServer(e)
}

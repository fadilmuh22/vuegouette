package internal

import (
	"context"
	"fmt"
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

		err := e.Start(fmt.Sprintf(":%d", viper.GetInt("PORT")))

		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server, ", err)
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

	postgreDB := db.Connect()
	redisClient := db.NewRedisClient()

	e.Use(middleware.DBMiddleware(postgreDB))
	e.Use(middleware.RedisMiddleware(redisClient))
	e.Use(middleware.TransformErrorResponse)

	// Validator
	var err error
	e.Validator, err = util.NewCustomValidator()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Routes
	handler.NewApiHandlers(e, postgreDB, redisClient)

	// Start server
	runServer(e)
}

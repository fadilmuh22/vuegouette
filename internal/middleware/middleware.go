package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/util"
)

func DBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(util.DBContextKey, db)
			return next(c)
		}
	}
}

func RedisMiddleware(redisClient *db.RedisClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(util.RedisContextKey, redisClient)
			return next(c)
		}
	}
}

func TransformErrorResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				return util.SendResponse(c, he.Code, false, he.Message.(string), nil)
			}

			return util.SendResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
		}

		return nil
	}
}

package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/util"
)

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

func Auth() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(util.Claims)
		},
		SigningKey:  []byte(util.GetJWTSecret()),
		TokenLookup: "header:Authorization",
	})
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*util.Claims)

		if !claims.IsAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

func DBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(util.DBContextKey, db)
			return next(c)
		}
	}
}

func ProductAuthor(productAuthor string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// user := c.Get("user").(*jwt.Token)
			// claims := user.Claims.(*util.Claims)
			db := c.Get(util.DBContextKey).(*gorm.DB)

			productID := c.Param("id")
			var product model.Product
			result := db.Where("id = ?", productID).First(&product)

			if result.Error != nil {
				return echo.NewHTTPError(http.StatusNotFound, "Product not found")
			}

			// if claims.ID != product.UserID {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			// }

			return next(c)
		}
	}
}

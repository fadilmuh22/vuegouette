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

var (
	Auth = echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(util.Claims)
		},
		SigningKey:  []byte(util.GetJWTSecret()),
		TokenLookup: "header:Authorization",
		ContextKey:  util.JWTContextKey,
		SuccessHandler: func(c echo.Context) {
			token := c.Get(util.JWTContextKey).(*jwt.Token)
			claims := token.Claims.(*util.Claims)

			c.Set(util.AuthContextKey, claims)
		},
	})
)

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Get(util.JWTContextKey).(*util.Claims)

		if !auth.IsAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

func ProductAuthor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Get(util.JWTContextKey).(*util.Claims)
		db := c.Get(util.DBContextKey).(*gorm.DB)

		productID := c.Param("id")
		var product model.Product
		result := db.Where("id = ?", productID).First(&product)

		if result.Error != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}

		if !auth.IsAdmin && auth.ID != product.UserID {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/fadilmuh22/restskuy/internal/util"
)

func Auth() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(util.Claims)
		},
		SigningKey:  []byte(util.GetJWTSecret()),
		TokenLookup: "header:Authorization",
		ContextKey:  util.JWTContextKey,
		SuccessHandler: func(c echo.Context) {
			user := c.Get(util.JWTContextKey).(*jwt.Token)
			claims := user.Claims.(*util.Claims)

			c.Set(util.AuthContextKey, claims)
		},
	})
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Get(util.AuthContextKey).(*util.Claims)

		if !auth.IsAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

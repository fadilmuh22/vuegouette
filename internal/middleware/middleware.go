package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

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

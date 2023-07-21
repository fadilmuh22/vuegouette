package util

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/fadilmuh22/restskuy/internal/model"
)

const (
	DBContextKey   = "__db"
	JWTContextKey  = "__jwt"
	AuthContextKey = "__auth"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password, err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func SendResponse(c echo.Context, status int, success bool, message string, data interface{}) error {
	return c.JSON(status, &model.BasicResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

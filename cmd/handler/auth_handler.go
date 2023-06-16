package handler

import (
	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/fadilmuh22/restskuy/cmd/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type authHandler struct {
	service service.AuthService
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseData struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (h authHandler) register(c echo.Context) error {
	var user model.User
	c.Bind(&user)

	user, err := h.service.Register(user)
	if err != nil {
		return err
	}

	return SendResponse(c, 200, true, "Success register", user)
}

func (h authHandler) login(c echo.Context) error {
	var loginRequestBody LoginRequestBody
	c.Bind(&loginRequestBody)

	user, err := h.service.Login(loginRequestBody.Email, loginRequestBody.Password)
	if err != nil {
		return err
	}

	// generate token
	accessToken, err := service.GenerateTokensAndSetCookies(&user, c)
	if err != nil {
		return err
	}

	return SendResponse(c, 200, true, "Success login", LoginResponseData{
		Token: accessToken,
		User:  user,
	})
}

func (h authHandler) HandleRoutes(g *echo.Group) {
	auth := g.Group("/auth")
	{
		auth.POST("/login", h.login)
		auth.POST("/register", h.register)
	}
}

func NewAuthHandler(db *gorm.DB) Handler {
	return authHandler{
		service: service.NewAuthService(db),
	}
}

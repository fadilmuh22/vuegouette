package handler

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/service"
	"github.com/fadilmuh22/restskuy/internal/util"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(db *gorm.DB) Handler {
	return authHandler{
		service: service.NewAuthService(db),
	}
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
		return echo.NewHTTPError(401, err)
	}

	return util.SendResponse(c, 200, true, "Success register", user)
}

func (h authHandler) login(c echo.Context) error {
	var loginRequestBody LoginRequestBody
	c.Bind(&loginRequestBody)

	user, err := h.service.Login(loginRequestBody.Email, loginRequestBody.Password)
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	// generate token
	accessToken, err := util.GenerateAccessToken(&user, c)
	if err != nil {
		return err
	}

	return util.SendResponse(c, 200, true, "Success login", LoginResponseData{
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

package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"

	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/middleware"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/service"
	"github.com/fadilmuh22/restskuy/internal/util"
)

type userHandler struct {
	service service.AuthService
}

func NewUserHandler(db *gorm.DB) Handler {
	return userHandler{
		service: service.NewAuthService(db),
	}
}

func (h userHandler) getAllUser(c echo.Context) error {
	var users []model.User

	users, err := h.service.FindAll()
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success get all user", users)
}

// get user by id
func (h userHandler) getUser(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := h.service.FindById(id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return util.SendResponse(c, http.StatusOK, true, "Success get user", user)
}

func (h userHandler) createUser(c echo.Context) error {
	var user model.User
	c.Bind(&user)

	if err := c.Validate(user); err != nil {
		return err
	}

	user, err := h.service.Create(user)
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success create user", user)
}

func (h userHandler) updateUser(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := h.service.FindById(id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	c.Bind(&user)
	user.ID = id

	if err := c.Validate(user); err != nil {
		return err
	}

	user.Password, err = util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user, err = h.service.Update(user)
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success update user", user)
}

func (h userHandler) deleteUser(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := h.service.Delete(model.User{ID: id})
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success delete user", user)
}

func (h userHandler) updateUserProfile(c echo.Context) error {
	auth := c.Get(util.AuthContextKey).(*util.Claims)

	var video model.TikTokItem
	c.Bind(&video)

	profile, err := h.service.UpdateUserProfileInterests(auth.User.ID, video)
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Scucess update profile", profile)

}

func (h userHandler) createGuestUser(c echo.Context) error {
	createUser := model.User{
		Name:     fmt.Sprintf("guest-%s", uuid.NewV4().String()),
		Email:    fmt.Sprintf("guest-%s@mail.com", uuid.NewV4().String()),
		Password: uuid.NewV4().String(),
		IsGuest:  true,
	}

	user, err := h.service.Create(createUser)
	if err != nil {
		return err
	}

	newProfile := model.UserProfile{
		UserID: user.ID,
	}

	_, err = h.service.CreateUserProfile(newProfile)
	if err != nil {
		return err
	}

	accessToken, err := util.GenerateAccessToken(&user, c)
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success create guest user", LoginResponseData{
		Token: accessToken,
		User:  user,
	})
}

func (h userHandler) HandleRoutes(g *echo.Group) {
	user := g.Group("/user")
	{
		user.GET("", h.getAllUser)
		user.GET("/:id", h.getUser)
		user.POST("", h.createUser, middleware.Auth(), middleware.Admin)
		user.PUT("/:id", h.updateUser, middleware.Auth(), middleware.Admin)
		user.DELETE("/:id", h.deleteUser, middleware.Auth(), middleware.Admin)

		user.POST("/guest", h.createGuestUser)
		user.PUT("/profile", h.updateUserProfile, middleware.Auth())
	}
}

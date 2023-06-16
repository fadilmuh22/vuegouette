package handler

import (
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userHandler struct {
	db *gorm.DB
}

func (h userHandler) getAllUser(c echo.Context) error {
	var users []model.User

	result := h.db.Find(&users)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success get all user", users)
}

// get user by id
func (h userHandler) getUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	user := model.User{UUID: id}

	result := h.db.First(&user)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success get user", user)
}

func (h userHandler) createUser(c echo.Context) error {
	var user model.User
	c.Bind(&user)

	result := h.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success create user", user)
}

func (h userHandler) updateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	user := model.User{UUID: id}

	result := h.db.First(&user)
	if result.Error != nil {
		return result.Error
	}

	c.Bind(&user)
	user.UUID = id
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result = h.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success update user", user)
}

func (h userHandler) deleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	user := model.User{UUID: id}

	result := h.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success delete user", nil)
}

func (h userHandler) HandleRoutes(g *echo.Group) {
	user := g.Group("/user")
	{
		user.GET("", h.getAllUser)
		user.GET("/:id", h.getUser)
		user.POST("", h.createUser)
		user.PUT("/:id", h.updateUser)
		user.DELETE("/:id", h.deleteUser)
	}
}

func NewUserHandler(db *gorm.DB) Handler {
	return userHandler{
		db: db,
	}
}

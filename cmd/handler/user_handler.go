package handler

import (
	"database/sql"
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/fadilmuh22/restskuy/cmd/service"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	service service.UserService
}

func (h userHandler) getAllUser(c echo.Context) error {
	users, err := h.service.GetAllUser()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success get all user",
		Data:    users,
	})
}

// get user by id
func (h userHandler) getUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success get user",
		Data:    user,
	})
}

func (h userHandler) createUser(c echo.Context) error {
	var user model.User
	c.Bind(&user)

	user, err := h.service.CreateUser(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success create user",
		Data:    user,
	})
}

func (h userHandler) updateUser(c echo.Context) error {
	var user model.User
	c.Bind(&user)

	id := c.Param("id")
	user, err := h.service.UpdateUser(id, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success update user",
		Data:    user,
	})
}

func (h userHandler) deleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.service.DeleteUser(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success delete user",
		Data:    nil,
	})
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

func NewUserHandler(con *sql.DB) Handler {
	return userHandler{
		service: service.UserService{
			Con: con,
		},
	}
}

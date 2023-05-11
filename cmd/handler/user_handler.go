package handler

import (
	"database/sql"
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/fadilmuh22/restskuy/cmd/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
}

func (h UserHandler) getAllUser(c echo.Context) error {
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
func (h UserHandler) getUser(c echo.Context) error {
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

func (h UserHandler) createUser(c echo.Context) error {
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

func (h UserHandler) updateUser(c echo.Context) error {
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

func (h UserHandler) deleteUser(c echo.Context) error {
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

func (h UserHandler) HandleRoutes(g *echo.Group) {
	user := g.Group("/user")
	{
		user.GET("", h.getAllUser)
		user.GET("/:id", h.getUser)
		user.POST("", h.createUser)
		user.PUT("/:id", h.updateUser)
		user.DELETE("/:id", h.deleteUser)
	}
}

func NewUserHandler(con *sql.DB) UserHandler {
	return UserHandler{
		service: service.UserService{
			Con: con,
		},
	}
}

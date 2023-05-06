package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getAllUser(c echo.Context) error {
	users, err := GetAllUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// get user by id
func getUser(c echo.Context) error {
	id := c.Param("id")
	user, err := GetUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func createUser(c echo.Context) error {
	var user User
	c.Bind(&user)

	user, err := CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func updateUser(c echo.Context) error {
	var user User
	c.Bind(&user)

	id := c.Param("id")
	user, err := UpdateUser(id, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	err := DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "User deleted")
}

func HandleRoutes(g *echo.Group) {
	user := g.Group("/user")
	{
		user.GET("", getAllUser)
		user.GET("/:id", getUser)
		user.POST("", createUser)
		user.PUT("/:id", updateUser)
		user.DELETE("/:id", deleteUser)
	}
}

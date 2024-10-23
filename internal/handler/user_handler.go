package handler

import (
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
	db := c.Get(util.DBContextKey).(*gorm.DB)

	var video model.TikTokItem
	c.Bind(&video)

	profile, err := h.service.GetUserProfile(auth.User.ID.String())
	if err != nil {
		return err
	}

	var titleAndTags []string
	titleAndTags = append(titleAndTags, util.TokenizeString(video.VideoTitle)...)
	titleAndTags = append(titleAndTags, video.Tags...)

	// Map to track changes or new interests
	updatedInterests := make(map[string]*model.Interest)

	// Process the title and tags to update the interests
	for _, term := range titleAndTags {
		found := false
		for i := range profile.Interests {
			if profile.Interests[i].Term == term {
				// Increment the interest score
				profile.Interests[i].WeightedScore++
				updatedInterests[term] = &profile.Interests[i]
				found = true
				break
			}
		}
		// If term not found in profile, add as a new interest
		if !found {
			newInterest := model.Interest{
				Term:          term,
				WeightedScore: 1,
				UserProfileID: profile.ID,
			}
			// Add the new interest to the map and profile
			updatedInterests[term] = &newInterest
			profile.Interests = append(profile.Interests, newInterest)
		}
	}

	// Update the user profile with the modified interests in the database
	if err := db.Transaction(func(tx *gorm.DB) error {
		// Update the profile itself
		if err := tx.Save(&profile).Error; err != nil {
			return err
		}

		// Update or insert interests
		for _, interest := range updatedInterests {
			if err := tx.Save(interest).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return util.SendResponse(c, http.StatusInternalServerError, false, "Failed to update profile in DB", nil)
	}

	return util.SendResponse(c, http.StatusOK, true, "Scucess update profile", profile)

}

func (h userHandler) createGuestUser(c echo.Context) error {
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

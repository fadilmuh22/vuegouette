package handler

import (
	"cmp"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/middleware"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/service"
	"github.com/fadilmuh22/restskuy/internal/util"
)

type videoHandler struct {
	videoService service.VideoService
	userService  service.UserService
}

func NewVideoHandler(db *gorm.DB, redisClient *db.RedisClient) Handler {
	return videoHandler{
		videoService: service.NewVideoService(db, redisClient),
		userService:  service.NewUserService(db),
	}
}

func (h videoHandler) getAllTikTokVideos(c echo.Context) error {
	keyword := c.QueryParam("keyword")

	videos, cached, err := h.videoService.FetchTikTokVideosWithCache(keyword)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Fetched from %s", util.IfThenElse(cached, "cache", "API"))

	return util.SendResponse(c, http.StatusOK, true, message, videos)
}

func (h videoHandler) getPersonalizedTiktokVideos(c echo.Context) error {
	var auth *util.Claims
	var keywords []string

	if c.Get(util.AuthContextKey) != nil {
		auth = c.Get(util.AuthContextKey).(*util.Claims)

		profile, err := h.userService.GetUserProfile(auth.User.ID.String())
		if err == nil {
			if profile.Interests != nil && len(profile.Interests) > 0 {
				interestRanked := profile.Interests

				slices.SortFunc(interestRanked, func(a, b model.Interest) int {
					return cmp.Compare(a.WeightedScore, b.WeightedScore)
				})

				for _, interest := range interestRanked[0:4] {
					keywords = append(keywords, interest.Term)
				}
			}
		}
	}

	if len(keywords) == 0 {
		keywords = append(keywords, "trending")
	}

	videos, cached, err := h.videoService.FetchTikTokVideosWithCache(strings.Join(keywords, "+"))
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Fetched from %s", util.IfThenElse(cached, "cache", "API"))

	return util.SendResponse(c, http.StatusOK, true, message, videos)
}

func fetchVideoLink(c echo.Context) error {
	videoURL := c.QueryParam("videoUrl") // Get videoUrl as query param

	// Prepare the request body for the external API
	requestBody := `{"language_id": 1, "query": "` + videoURL + `"}`

	// Send the POST request to the external API
	req, err := http.NewRequest("POST", "https://ttsave.app/download", io.NopCloser(strings.NewReader(requestBody)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create request"})
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch data"})
	}
	defer resp.Body.Close()

	// Read the response from the external API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read response"})
	}

	// Return the API response as-is
	return util.SendResponse(c, http.StatusOK, true, "Fetched from API", string(body))
}

func (h videoHandler) HandleRoutes(g *echo.Group) {
	video := g.Group("/video")
	{
		video.GET("", h.getAllTikTokVideos)
		video.GET("/personalized", h.getPersonalizedTiktokVideos, middleware.Guest())
		video.GET("/fetch-video", fetchVideoLink)
	}
}

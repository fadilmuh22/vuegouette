package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/middleware"
	"github.com/fadilmuh22/restskuy/internal/service"
	"github.com/fadilmuh22/restskuy/internal/util"
)

type videoHandler struct {
	videoService service.VideoService
	userService  service.UserService
}

func NewVideoHandler(db *gorm.DB, redisClient *db.RedisClient) Handler {
	return &videoHandler{
		videoService: service.NewVideoService(db, redisClient),
		userService:  service.NewUserService(db),
	}
}

func (h *videoHandler) searchTikTokVideos(c echo.Context) error {
	var userID string

	keyword := c.QueryParam("keyword")
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")

	pageNum, _ := strconv.Atoi(util.IfThenElse(page == "", "1", page).(string))
	pageSizeNum, _ := strconv.Atoi(util.IfThenElse(pageSize == "", "10", pageSize).(string))

	auth := util.TryGetAuth(c)

	if auth == nil {
		userID = fmt.Sprint("guest-", uuid.NewV4().String())
	} else {
		userID = auth.User.ID.String()
	}

	videos, cached, err := h.videoService.FetchTikTokVideosWithCache(userID, keyword, pageNum, pageSizeNum)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Fetched from %s", util.IfThenElse(cached, "cache", "API"))

	return util.SendResponse(c, http.StatusOK, true, message, videos)
}

func (h *videoHandler) getPersonalizedTiktokVideos(c echo.Context) error {
	var userID string

	keyword := c.QueryParam("keyword")
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")

	pageNum, _ := strconv.Atoi(util.IfThenElse(page == "", "1", page).(string))
	pageSizeNum, _ := strconv.Atoi(util.IfThenElse(pageSize == "", "10", pageSize).(string))

	auth := util.TryGetAuth(c)

	if auth == nil {
		userID = fmt.Sprint("guest-", uuid.NewV4().String())
	} else {
		userID = auth.User.ID.String()
	}

	if keyword == "" {
		keyword = "trending"
	}

	videos, cached, err := h.videoService.FetchTikTokVideosWithCache(userID, keyword, pageNum, pageSizeNum)
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch data")
	}
	defer resp.Body.Close()

	// Read the response from the external API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read response")
	}

	// Return the API response as-is
	return util.SendResponse(c, http.StatusOK, true, "Fetched from API", string(body))
}

func (h *videoHandler) HandleRoutes(g *echo.Group) {
	video := g.Group("/video")
	{
		video.GET("", h.searchTikTokVideos, middleware.Guest())
		video.GET("/personalized", h.getPersonalizedTiktokVideos, middleware.Guest())
		video.GET("/fetch-video", fetchVideoLink, middleware.Guest())
	}
}

package handler

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/db"
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
	auth := c.Get(util.AuthContextKey).(*util.Claims)

	profile, err := h.userService.GetUserProfile(auth.User.ID.String())
	if err != nil {
		return err
	}

	var keywords []string
	if profile.Interests != nil {
		interestRanked := profile.Interests

		slices.SortFunc(interestRanked, func(a, b model.Interest) int {
			return cmp.Compare(a.WeightedScore, b.WeightedScore)
		})

		for _, interest := range interestRanked[0:4] {
			keywords = append(keywords, interest.Term)
		}

	}

	videos, cached, err := h.videoService.FetchTikTokVideosWithCache(strings.Join(keywords, "+"))
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Fetched from %s", util.IfThenElse(cached, "cache", "API"))

	return util.SendResponse(c, http.StatusOK, true, message, videos)
}

func (h videoHandler) asyncFetchTikTokVideos(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	resultChan := make(chan []model.TikTokItem)

	go func() {
		videos, cached, err := h.videoService.FetchTikTokVideosWithCache(keyword)
		if err != nil {
			log.Println("Error fetching videos:", err, cached)
			resultChan <- []model.TikTokItem{}
		} else {
			resultChan <- videos
		}
	}()

	return util.SendResponse(c, http.StatusOK, true, "Fetched videos", <-resultChan)

}

func (h videoHandler) HandleRoutes(g *echo.Group) {
	video := g.Group("/video")
	{
		video.GET("", h.getAllTikTokVideos)
		video.GET("/async", h.asyncFetchTikTokVideos)
	}
}

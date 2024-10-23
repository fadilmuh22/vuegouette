package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/service"
	"github.com/fadilmuh22/restskuy/internal/util"
)

type videoHandler struct {
	service service.VideoService
}

func NewVideoHandler(db *gorm.DB, redisClient *db.RedisClient) Handler {
	return videoHandler{
		service: service.NewVideoService(db, redisClient),
	}
}

func (h videoHandler) getAllTikTokVideos(c echo.Context) error {
	keyword := c.QueryParam("keyword")

	videos, cached, err := h.service.FetchTikTokVideosWithCache(keyword)
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
		videos, cached, err := h.service.FetchTikTokVideosWithCache(keyword)
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

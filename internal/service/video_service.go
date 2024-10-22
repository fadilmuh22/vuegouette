package service

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/util"
	"gorm.io/gorm"
)

type VideoService struct {
	db *gorm.DB
}

func NewVideoService(db *gorm.DB) VideoService {
	return VideoService{db: db}
}

func (s VideoService) fetchTikTokVideos(keyword string) ([]string, error) {
	url := fmt.Sprintf("https://www.tiktok.com/search?q=%s", keyword)

	// Create HTTP request with user-agent to mimic a browser request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract video URLs (this part will depend on TikTok's page structure)
	var videoURLs []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists && util.IsValidTikTokURL(url) {
			videoURLs = append(videoURLs, url)
		}
	})

	// Save the video URLs to the database
	videoKeyword := model.VideoKeyword{
		Keyword:   keyword,
		VideoUrls: videoURLs,
	}
	err = s.db.Create(&videoKeyword).Error
	if err != nil {
		return nil, err
	}

	return videoURLs, nil
}

func (s VideoService) FetchTikTokVideosWithCache(keyword string) ([]string, bool, error) {
	var videoKeyword model.VideoKeyword

	err := s.db.Where("keyword = ?", keyword).First(&videoKeyword).Error
	if err == nil {
		if len(videoKeyword.VideoUrls) > 0 {
			return videoKeyword.VideoUrls, true, nil
		}
	}

	videos, err := s.fetchTikTokVideos(keyword)
	if err == nil {
		return videos, false, nil
	}
	return videos, false, err
}

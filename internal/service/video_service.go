package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/model"

	"gorm.io/gorm"
)

type VideoService struct {
	db          *gorm.DB
	redisClient *db.RedisClient
}

func NewVideoService(db *gorm.DB, resdisClient *db.RedisClient) VideoService {
	return VideoService{db: db, redisClient: resdisClient}
}

func (s VideoService) fetchTikTokVideos(keyword string) ([]model.TikTokItem, error) {
	log.Printf("Fetching keyword: %s\n", keyword)
	url := fmt.Sprintf("https://www.tiktok.com/search?q=%s", keyword)

	optsAlloc := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.Flag(`headless`, true),
		chromedp.DisableGPU,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), optsAlloc...)
	defer cancel()

	var optsTask []chromedp.ContextOption
	optsTask = append(optsTask, chromedp.WithDebugf(log.Printf))

	// Create a new context
	taskCtx, cancel := chromedp.NewContext(allocCtx, optsTask...)
	defer cancel()

	var tiktokItems []model.TikTokItem

	// Run the tasks
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#main-content-general_search", chromedp.ByID), // Wait for specific element to be visible
		chromedp.Sleep(3*time.Second),                                       // Wait for 2 seconds
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Evaluate JavaScript to extract TikTok video information
			return chromedp.Evaluate(`Array.from(document.querySelectorAll("div[class*='DivItemContainerForSearch']")).map(item => {
				return {
					video_url: item.querySelector("a[class*='AVideoContainer']")?.href,
					avatar_url: item.querySelector("span[class*='SpanAvatarContainer'] img")?.src,
					user_name: item.querySelector("p[class*='PUniqueId']")?.innerText,
					video_title: item.querySelector("div[data-e2e='search-card-video-caption']")?.
						querySelector("span[class*='SpanText']")?.innerText,
					video_count: item.querySelector("strong[class*='StrongVideoCount']")?.innerText,
					tags: Array.from(item.querySelectorAll("a[data-e2e='search-common-link']")).map(tag => tag.href),
				};
			});`, &tiktokItems).Do(ctx)
		}),
	)

	if err != nil {
		return nil, err
	}

	return tiktokItems, nil
}

func (s VideoService) FetchTikTokVideosWithCache(keyword string) ([]model.TikTokItem, bool, error) {
	cachedVideos, err := s.redisClient.Get(keyword)
	if err == nil {
		var cachedItems []model.TikTokItem
		err := json.Unmarshal([]byte(cachedVideos), &cachedItems)
		if err == nil && len(cachedItems) > 0 {
			log.Println("Returning cached videos")
			return cachedItems, true, nil
		}
	}

	videos, err := s.fetchTikTokVideos(keyword)
	if err == nil {
		// Cache the results in Redis
		cachedData, err := json.Marshal(videos)
		if err == nil {
			err = s.redisClient.Set(keyword, cachedData, 10*time.Minute) // Set expiration as needed
			if err != nil {
				log.Println("Error setting cache:", err)
			}
		}

		return videos, false, nil
	}
	return videos, false, err
}

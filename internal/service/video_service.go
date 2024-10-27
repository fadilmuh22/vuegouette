package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"

	"gorm.io/gorm"
)

type userContext struct {
	ctx      context.Context
	cancel   context.CancelFunc
	lastUsed time.Time
}

type VideoService struct {
	db           *gorm.DB
	redisClient  *db.RedisClient
	userContexts map[string]*userContext // Maps userID to Chromedp context
	mutex        sync.Mutex              // Protects userContexts map
	expireAfter  time.Duration           // Duration to wait before expiring a context
}

func NewVideoService(db *gorm.DB, redisClient *db.RedisClient) VideoService {
	return VideoService{
		db:           db,
		redisClient:  redisClient,
		userContexts: make(map[string]*userContext),
		expireAfter:  10 * time.Minute,
	}
}

func (s *VideoService) getOrCreateUserContext(userID string) (context.Context, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Infof("Reusing existing context for user %s", userID, s.userContexts[userID])
	// Check if context already exists for the user
	if uCtx, exists := s.userContexts[userID]; exists {
		return uCtx.ctx, nil
	}

	// Create new Chromedp context for the user
	optsAlloc := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.Flag("headless", true),
		chromedp.DisableGPU,
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), optsAlloc...)
	taskCtx, cancel := chromedp.NewContext(allocCtx)

	// Store the context with expiration handling
	s.userContexts[userID] = &userContext{
		ctx:      taskCtx,
		cancel:   cancel,
		lastUsed: time.Now(),
	}
	return taskCtx, nil
}

func (s *VideoService) fetchTikTokVideos(userID, keyword string) ([]model.TikTokItem, error) {
	url := fmt.Sprintf("https://www.tiktok.com/search?q=%s", keyword)

	// Get or create user-specific Chromedp context
	taskCtx, err := s.getOrCreateUserContext(userID)
	if err != nil {
		return nil, err
	}

	var tiktokItems []model.TikTokItem

	// Run the tasks
	err = chromedp.Run(taskCtx,
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
					tags: Array.from(item.querySelectorAll("a[data-e2e='search-common-link']"))
						.map(tag => tag.href)
						.filter(tag => tag.includes("/tag"))
						.map(tag => tag.split("/").pop())
				};
			});`, &tiktokItems).Do(ctx)
		}),
	)

	if err != nil {
		return nil, err
	}

	return tiktokItems, nil
}

func (s *VideoService) FetchTikTokVideosWithCache(userID, keyword string) ([]model.TikTokItem, bool, error) {
	log.Info("Fetching keywords: ", keyword)
	cachedVideos, err := s.redisClient.Get(keyword)
	if err == nil {
		var cachedItems []model.TikTokItem
		err := json.Unmarshal([]byte(cachedVideos), &cachedItems)
		if err == nil && len(cachedItems) > 0 {
			log.Info("Returning cached videos")
			return cachedItems, true, nil
		}
	}

	videos, err := s.fetchTikTokVideos(userID, keyword)
	if err == nil {
		if len(videos) == 0 {
			return videos, false, model.NewErrorMessage("No video available at the moment")
		}

		// Add ID to each video item
		for i := 0; i < len(videos); i++ {
			videos[i].ID = uuid.NewV4()
		}

		// Cache the results in Redis
		cachedData, err := json.Marshal(videos)
		if err == nil {
			err = s.redisClient.Set(keyword, cachedData, 10*time.Minute) // Set expiration as needed
			if err != nil {
				log.Debug("Error setting cache:", err)
			}
		}

		return videos, false, nil
	}
	return videos, false, err
}

func (s *VideoService) cleanupExpiredContexts() {
	ticker := time.NewTicker(5 * time.Minute) // Run cleanup every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		for userID, uCtx := range s.userContexts {
			if time.Since(uCtx.lastUsed) > s.expireAfter {
				uCtx.cancel() // Close the context
				delete(s.userContexts, userID)
				log.Infof("Expired Chromedp context for user: %s", userID)
			}
		}
		s.mutex.Unlock()
	}
}

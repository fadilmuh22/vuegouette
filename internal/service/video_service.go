package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/fadilmuh22/restskuy/internal/db"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/labstack/gommon/log"

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

func (s *VideoService) getOrCreateUserContext(userID string) (context.Context, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Infof("Reusing existing context for user %s", userID, s.userContexts[userID])
	// Check if context already exists for the user
	if uCtx, exists := s.userContexts[userID]; exists {
		return uCtx.ctx, true, nil
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
	return taskCtx, false, nil
}

func (s *VideoService) fetchTikTokVideos(userID, keyword string, page, pageSize int) ([]model.TikTokItem, error) {
	url := fmt.Sprintf("https://www.tiktok.com/search?q=%s", keyword)

	// Get or create user-specific Chromedp context
	taskCtx, existed, err := s.getOrCreateUserContext(fmt.Sprint(userID, '-', keyword))
	if err != nil {
		return nil, err
	}

	// Slice to hold the scraped video items
	var tiktokItems []model.TikTokItem

	// Navigate to the TikTok search page
	if !existed {
		err = chromedp.Run(taskCtx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("#main-content-general_search", chromedp.ByID),
			chromedp.Sleep(3*time.Second),
		)
	}

	if err != nil {
		log.Errorf("Error navigating to URL: %v", err)
		return nil, err
	}

	// Scroll and collect videos based on the requested page
	startIndex := page * pageSize
	targetIndex := startIndex + pageSize
	currentIndex := 0

	for currentIndex < targetIndex {
		// Scraping logic to extract video data
		var items []model.TikTokItem
		err := chromedp.Run(taskCtx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				return s.extractTikTokVideos(ctx, currentIndex, &items)
			}),
		)

		if err != nil {
			log.Errorf("Error extracting videos: %v", err)
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			if currentIndex >= startIndex && currentIndex < targetIndex {
				tiktokItems = append(tiktokItems, item)
			}
			currentIndex++
			if currentIndex >= targetIndex {
				break
			}
		}

		// Scroll to load more videos if needed
		if currentIndex < targetIndex {
			err = chromedp.Run(taskCtx,
				chromedp.KeyEvent(kb.End),
				chromedp.Sleep(3*time.Second),
			)
			if err != nil {
				log.Errorf("Error scrolling page: %v", err)
				return nil, err
			}
		}
	}

	return tiktokItems, nil
}

func (s *VideoService) extractTikTokVideos(ctx context.Context, currentIndex int, tiktokItems *[]model.TikTokItem) error {
	return chromedp.Evaluate(fmt.Sprintf(`Array.from(document.querySelectorAll("div[class*='DivItemContainerForSearch']")).slice(%d).map(item => {
        return {
            video_url: item.querySelector("a[class*='AVideoContainer']")?.href,
            avatar_url: item.querySelector("span[class*='SpanAvatarContainer'] img")?.src,
            user_name: item.querySelector("p[class*='PUniqueId']")?.innerText,
            video_title: item.querySelector("div[data-e2e='search-card-video-caption']")?.
                querySelector("span[class*='SpanText']")?.innerText,
            video_image_url: item.querySelector("picture > img")?.src,
            video_count: item.querySelector("strong[class*='StrongVideoCount']")?.innerText,
            tags: Array.from(item.querySelectorAll("a[data-e2e='search-common-link']"))
                .map(tag => tag.href)
                .filter(tag => tag.includes("/tag"))
                .map(tag => tag.split("/").pop())
        };
    });`, currentIndex), tiktokItems).Do(ctx)
}

func (s *VideoService) FetchTikTokVideosWithCache(userID, keyword string, page, pageSize int) ([]model.TikTokItem, bool, error) {
	log.Infof("Fetching videos for keyword: %s, page: %d, and pageSize: %d", keyword, page, pageSize)

	// Create a Redis key for tracking seen videos
	seenVideosKey := fmt.Sprintf("seen_videos:%s:%s", userID, keyword)

	videos, err := s.fetchTikTokVideos(userID, keyword, page, pageSize)
	if err != nil {
		return nil, false, err
	}

	// Filter out duplicate videos based on seen video IDs in Redis
	uniqueVideos := make([]model.TikTokItem, 0, len(videos))
	for _, video := range videos {
		// Generate unique identifier for each video (could use video URL or video ID)
		videoID := video.ID.String()
		isSeen, err := s.redisClient.SIsMember(seenVideosKey, videoID).Result()
		if err != nil {
			log.Error("Error checking Redis set:", err)
			continue
		}

		// Add video to the result if it hasn't been seen yet
		if !isSeen {
			uniqueVideos = append(uniqueVideos, video)

			// Mark video as seen by adding its ID to the Redis set
			if err := s.redisClient.SAdd(seenVideosKey, videoID).Err(); err != nil {
				log.Error("Error adding video ID to Redis set:", err)
			}
		}
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

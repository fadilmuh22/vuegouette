package util

import (
	"fmt"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/fadilmuh22/restskuy/internal/model"
)

const (
	DBContextKey    = "__db"
	RedisContextKey = "__redis"
	JWTContextKey   = "__user"
	AuthContextKey  = "__auth"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password, err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func SendResponse(c echo.Context, status int, success bool, message string, data interface{}) error {
	return c.JSON(status, &model.BasicResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func TokenizeString(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func DictContainWithSubstring(dict []string, text string) bool {
	for i := range dict {
		if strings.Contains(text, dict[i]) || text == dict[i] {
			return true
		}
	}
	return false
}

func TokenizeTikokItem(video model.TikTokItem) []string {
	var videoTags []string
	videoTags = append(videoTags, Map(video.Tags, func(tag string) string {
		return strings.ToLower(tag)
	})...)

	// dicts
	dicts := []string{
		"tiktok",
		"fyp",
		"foryou",
		"foryoupage",
	}

	// sanitize title and tags
	var sanitizeVideoTags []string

	for i := range videoTags {
		if len(videoTags[i]) < 3 || DictContainWithSubstring(dicts, videoTags[i]) {
			continue
		}

		sanitizeVideoTags = append(sanitizeVideoTags, videoTags[i])
	}

	return sanitizeVideoTags
}

func UpdateInterestsWithSubstrings(profileID uuid.UUID, interests []model.Interest, titleAndTags []string) map[string]model.Interest {
	updatedInterests := make(map[string]model.Interest)

	// Loop through each term in title and tags
	for _, term := range titleAndTags {
		found := false

		// Check if the term (or a substring match) exists in the user's interests
		for i := range interests {
			existingInterest := strings.ToLower(interests[i].Term)

			// If the term matches or is a substring of an existing interest, increment score
			if existingInterest == term || strings.Contains(existingInterest, term) || strings.Contains(term, existingInterest) {
				interests[i].WeightedScore++
				updatedInterests[term] = interests[i]
				found = true
				break
			}
		}

		// If the term was not found as a substring match, add it as a new interest
		if !found {
			newInterest := model.Interest{
				Term:          term,
				WeightedScore: 1,
				UserProfileID: profileID,
			}
			updatedInterests[term] = newInterest
			interests = append(interests, newInterest)
		}
	}

	// sort interests by weighted score
	slices.SortFunc(interests, func(a, b model.Interest) int {
		return int(b.WeightedScore - a.WeightedScore)
	})

	// keep upper 5
	interests = interests[0:5]

	return updatedInterests
}

func TryGetAuth(c echo.Context) *Claims {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in getAuth:", r)
		}
	}()

	if claims, ok := c.Get(AuthContextKey).(*Claims); ok {
		return claims
	}

	// Return nil if the type assertion fails
	return nil
}

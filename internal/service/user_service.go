package service

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/util"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return UserService{db: db}
}

func (s UserService) FindAll() ([]model.User, error) {
	var users []model.User

	err := s.db.Model(&model.User{}).Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s UserService) FindById(id string) (model.User, error) {
	var user model.User

	err := s.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserService) FindByEmail(email string) (model.User, error) {
	var user model.User

	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserService) Create(user model.User) (model.User, error) {
	var err error
	user.Password, err = util.HashPassword(user.Password)
	if err != nil {
		return user, err
	}

	err = s.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserService) Update(user model.User) (model.User, error) {
	err := s.db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserService) Delete(user model.User) (model.User, error) {
	err := s.db.Delete(user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserService) GetUserProfile(userID string) (model.UserProfile, error) {
	var profile model.UserProfile

	err := s.db.Where("user_id = ?", userID).Preload("Interests", func(db *gorm.DB) *gorm.DB {
		return db.Order("interests.weighted_score DESC").Limit(5)
	}).First(&profile).Error
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s UserService) SaveUserProfile(profile *model.UserProfile) error {
	// Update the user profile in the database
	// Use Save or Update method depending on your use case
	// If the profile already exists, it will update the existing one
	// Otherwise, it will insert a new record
	err := s.db.Save(profile).Error
	if err != nil {
		return fmt.Errorf("could not save user profile: %w", err)
	}
	return nil
}


func (s UserService) GetUserProfileKeywords(userID string) []string {
	var keywords []string

	profile, err := s.GetUserProfile(userID)
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

	return keywords
}

func (s UserService) DeleteUserProfileKeywords(userID string, keywordsToDelete []string) ([]string, error) {
	var updatedKeywords []string

	// Fetch the user profile
	profile, err := s.GetUserProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("could not find user profile: %w", err)
	}

	// Loop through the list of keywords to delete
	for _, keyword := range keywordsToDelete {
		// Check if the interest exists for this user and delete it
		var interest model.Interest
		err := s.db.Where("user_profile_id = ? AND term = ?", profile.ID, keyword).First(&interest).Error
		if err == nil {
			// If the interest exists, delete it
			if deleteErr := s.db.Delete(&interest).Error; deleteErr != nil {
				return nil, fmt.Errorf("failed to delete interest: %w", deleteErr)
			}
		} else if err != gorm.ErrRecordNotFound {
			// If the error is not "record not found", return it
			return nil, fmt.Errorf("failed to check interest: %w", err)
		}
	}

	updatedKeywords = append(updatedKeywords, s.GetUserProfileKeywords(userID)...)

	return updatedKeywords, nil
}


func (s UserService) CreateUserProfile(profile model.UserProfile) (model.UserProfile, error) {
	err := s.db.Create(&profile).Error
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s UserService) UpdateUserProfile(profile model.UserProfile) (model.UserProfile, error) {
	err := s.db.Save(&profile).Error
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s UserService) UpdateUserProfileInterests(userID uuid.UUID, video model.TikTokItem) (model.UserProfile, error) {
	profile, err := s.GetUserProfile(userID.String())
	if err != nil {
		return profile, err
	}

	videoTags := util.TokenizeTikokItem(video)

	// Map to track changes or new interests
	updatedInterests := util.UpdateInterestsWithSubstrings(profile.ID, profile.Interests, videoTags)

	// Update the user profile with the modified interests in the database
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// Update or insert interests
		for _, interest := range updatedInterests {
			// check if term exist
			var existingInterest model.Interest
			tx.Where("term = ? AND user_profile_id = ?", interest.Term, profile.ID).Limit(1).Find(&existingInterest)

			if (model.Interest{} == existingInterest) {
				if err := tx.Create(&interest).Error; err != nil {
					return err
				}
				continue
			}

			existingInterest.WeightedScore = interest.WeightedScore
			if err := tx.Save(&existingInterest).Error; err != nil {
				return err
			}

		}

		return nil
	}); err != nil {
		return profile, model.NewErrorMessage("Failed to update user interests")
	}

	// delete all user interests below blow top 5 when ordered by weighted_score
	if err := s.db.Exec("DELETE FROM interests WHERE user_profile_id = ? AND id NOT IN (SELECT id FROM interests WHERE user_profile_id = ? ORDER BY weighted_score DESC LIMIT 5)", profile.ID, profile.ID).Error; err != nil {
		log.Error("Failed to delete user interests", err)
	}

	profile, err = s.GetUserProfile(userID.String())
	if err != nil {
		return profile, err
	}

	return profile, nil
}

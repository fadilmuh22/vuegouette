package service

import (
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/util"
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

func (s UserService) GetUserProfile(id string) (model.UserProfile, error) {
	var profile model.UserProfile

	err := s.db.Where("id = ?", id).First(&profile).Error
	if err != nil {
		return profile, err
	}

	return profile, nil
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

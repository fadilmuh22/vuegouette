package service

import (
	"github.com/fadilmuh22/restskuy/internal/model"
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

	result := s.db.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func (s UserService) FindById(id string) (model.User, error) {
	var user model.User

	result := s.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (s UserService) FindByEmail(email string) (model.User, error) {
	var user model.User

	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, nil
	}

	return user, nil
}

func (s UserService) Create(user model.User) (model.User, error) {
	result := s.db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (s UserService) Update(user model.User) (model.User, error) {
	result := s.db.Save(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (s UserService) Delete(user model.User) (model.User, error) {
	result := s.db.Delete(user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

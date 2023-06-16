package service

import (
	"github.com/fadilmuh22/restskuy/cmd/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return AuthService{db: db}
}

func (s AuthService) Register(user model.User) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashedPassword)

	result := s.db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (s AuthService) Login(email, password string) (model.User, error) {
	var user model.User

	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, model.NewErrorMessage("Email is not registered")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, model.NewErrorMessage("Password is wrong")
	}

	return user, nil
}

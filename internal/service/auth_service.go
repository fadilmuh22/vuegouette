package service

import (
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/util"
	"github.com/labstack/echo/v4"
)

type AuthService struct {
	UserService
}

func NewAuthService(db *gorm.DB) AuthService {
	return AuthService{
		UserService: NewUserService(db),
	}
}

func (s AuthService) Register(user model.User) (model.User, error) {
	user, err := s.UserService.Create(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s AuthService) Login(email, password string) (model.User, error) {
	var user model.User

	user, err := s.UserService.FindByEmail(email)
	if err != nil {
		return user, echo.NewHTTPError(401, model.NewErrorMessage("Email is not registered"))
	}

	err = util.ComparePassword(user.Password, password)
	if err != nil {
		return user, echo.NewHTTPError(401, model.NewErrorMessage("Password is incorrect"))
	}

	return user, nil
}

package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	Password string    `json:"password" gorm:"not null" validate:"required,min=6"`
	IsAdmin  bool      `json:"is_admin" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return
}

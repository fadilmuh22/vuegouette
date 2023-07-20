package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email" gorm:"uniqueIndex"`
	Password string    `json:"password" validate:"required,min=6" gorm:"not null"`
	IsAdmin  bool      `json:"is_admin" gorm:"default:false"`
	Products []Product `json:"products"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return
}

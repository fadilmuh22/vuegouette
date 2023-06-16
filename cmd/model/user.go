package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID     uuid.UUID `json:"uuid" gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	return
}

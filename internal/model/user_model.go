package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" gorm:"uniqueIndex" validate:"email"`
	Password string    `json:"password" gorm:"not null" validate:"min=6"`
	IsAdmin  bool      `json:"is_admin" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return
}

type UserProfile struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID  `gorm:"type:char(36);unique;not null"`             // Link to the user
	Interests []Interest `json:"interests" gorm:"foreignKey:UserProfileID"` // Interests calculated
}

func (u *UserProfile) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return
}

type Interest struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserProfileID uuid.UUID `gorm:"type:char(36);not null"`
	Term          string    `json:"term" gorm:"type:text"`
	WeightedScore float64   `json:"weighted_score"`
}

func (i *Interest) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewV4()
	return
}

func (i *Interest) ToMap() map[string]float64 {
	mapInterest := make(map[string]float64)
	mapInterest[i.Term] = i.WeightedScore
	return mapInterest
}

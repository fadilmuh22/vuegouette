package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type VideoKeyword struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Keyword   string    `json:"keyword" validate:"required"`
	VideoUrls []string  `json:"video_urls" gorm:"type:text"`
}

func (u *VideoKeyword) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return
}

type TikTokItem struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	AvatarURL  string    `json:"avatar_url" gorm:"type:text"`
	UserName   string    `json:"user_name" gorm:"type:text"`
	VideoTitle string    `json:"video_title" gorm:"type:text"`
	VideoCount string    `json:"video_count" gorm:"type:text"`
	VideoURL   string    `json:"video_url" gorm:"type:text"`
	Tags       []string  `json:"tags" gorm:"type:text"`
}

func (u *TikTokItem) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4()
	return
}
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

package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required;min=0"`
	Description *string   `json:"description" validate:"omitempty"`
	Stock       int       `json:"stock" validate:"min=0"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	User        User      `json:"user"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewV4()
	return
}

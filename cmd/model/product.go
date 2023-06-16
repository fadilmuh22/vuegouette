package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	UUID        uuid.UUID `json:"uuid" gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description *string   `json:"description"`
	Stock       int       `json:"stock"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.UUID = uuid.New()
	return
}

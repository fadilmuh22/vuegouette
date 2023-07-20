package service

import (
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/model"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return ProductService{db: db}
}

func (s ProductService) FindAll() ([]model.Product, error) {
	var products []model.Product

	err := s.db.Model(&model.Product{}).Preload("User").Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}

func (s ProductService) FindById(id string) (model.Product, error) {
	var product model.Product

	result := s.db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func (s ProductService) Create(product model.Product) (model.Product, error) {
	result := s.db.Create(&product)
	if result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func (s ProductService) Update(product model.Product) (model.Product, error) {
	result := s.db.Save(&product)
	if result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func (s ProductService) Delete(product model.Product) (model.Product, error) {
	result := s.db.Delete(&product)
	if result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

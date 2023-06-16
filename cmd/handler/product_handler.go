package handler

import (
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type productHandler struct {
	db *gorm.DB
}

// getProducts with query
func (h productHandler) getProducts(c echo.Context) error {
	var products []model.Product

	result := h.db.Find(&products)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success get all product", products)
}

func (h productHandler) getProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	product := model.Product{UUID: id}

	result := h.db.First(&product)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success get product", product)
}

func (h productHandler) createProduct(c echo.Context) error {
	var product model.Product
	c.Bind(&product)

	result := h.db.Create(&product)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success create product", product)
}

func (h productHandler) updateProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	product := model.Product{UUID: id}

	// find product by id
	result := h.db.First(&product)
	if result.Error != nil {
		return result.Error
	}

	c.Bind(&product)
	product.UUID = id

	result = h.db.Save(&product)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success update product", product)
}

func (h productHandler) deleteProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	product := model.Product{UUID: id}

	result := h.db.Delete(&product)
	if result.Error != nil {
		return result.Error
	}

	return SendResponse(c, http.StatusOK, true, "Success delete product", nil)
}

func (h productHandler) HandleRoutes(g *echo.Group) {
	product := g.Group("/product")
	{
		product.GET("", h.getProducts)
		product.POST("", h.createProduct)
		product.GET("/:id", h.getProduct)
		product.PUT("/:id", h.updateProduct)
		product.DELETE("/:id", h.deleteProduct)
	}
}

func NewProductHandler(db *gorm.DB) Handler {
	return productHandler{
		db: db,
	}
}

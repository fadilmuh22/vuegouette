package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/fadilmuh22/restskuy/internal/middleware"
	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/fadilmuh22/restskuy/internal/service"
	"github.com/fadilmuh22/restskuy/internal/util"
)

type productHandler struct {
	service service.ProductService
}

func NewProductHandler(db *gorm.DB) Handler {
	return productHandler{
		service: service.NewProductService(db),
	}
}

func (h productHandler) getProducts(c echo.Context) error {
	var products []model.Product

	products, err := h.service.FindAll()
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success get all product", products)
}

func (h productHandler) getProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	product, err := h.service.FindById(id.String())
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success get product", product)
}

func (h productHandler) createProduct(c echo.Context) error {
	var product model.Product
	c.Bind(&product)

	product, err := h.service.Create(product)
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success create product", product)
}

func (h productHandler) updateProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	product, err := h.service.FindById(id.String())
	if err != nil {
		return err
	}

	c.Bind(&product)
	product.UUID = id

	product, err = h.service.Update(product)
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success update product", product)
}

func (h productHandler) deleteProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	product, err := h.service.Delete(model.Product{UUID: id})
	if err != nil {
		return err
	}

	return util.SendResponse(c, http.StatusOK, true, "Success delete product", product)
}

func (h productHandler) HandleRoutes(g *echo.Group) {
	product := g.Group("/product")
	{
		product.GET("", h.getProducts)
		product.POST("", h.createProduct, middleware.Auth())
		product.GET("/:id", h.getProduct, middleware.Auth())
		product.PUT("/:id", h.updateProduct, middleware.Auth())
		product.DELETE("/:id", h.deleteProduct, middleware.Auth())
	}
}

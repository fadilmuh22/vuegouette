package handler

import (
	"database/sql"
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/fadilmuh22/restskuy/cmd/service"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	service service.ProductService
}

// getProducts with query
func (h productHandler) getProducts(c echo.Context) error {
	products, err := h.service.GetAllProduct()
	if err != nil {
		return err
	}

	return SendResponse(c, http.StatusOK, true, "Success get all product", products)
}

func (h productHandler) getProduct(c echo.Context) error {
	id := c.Param("id")
	product, err := h.service.GetProduct(id)
	if err != nil {
		return err
	}

	return SendResponse(c, http.StatusOK, true, "Success get product", product)
}

func (h productHandler) createProduct(c echo.Context) error {
	var product model.Product
	c.Bind(&product)

	product, err := h.service.CreateProduct(product)
	if err != nil {
		return err
	}

	return SendResponse(c, http.StatusOK, true, "Success create product", product)
}

func (h productHandler) updateProduct(c echo.Context) error {
	var product model.Product
	c.Bind(&product)

	id := c.Param("id")
	product, err := h.service.UpdateProduct(id, product)
	if err != nil {
		return err
	}

	return SendResponse(c, http.StatusOK, true, "Success update product", product)
}

func (h productHandler) deleteProduct(c echo.Context) error {
	id := c.Param("id")
	err := h.service.DeleteProduct(id)
	if err != nil {
		return err
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

func NewProductHandler(con *sql.DB) Handler {
	return productHandler{
		service: service.ProductService{
			Con: con,
		},
	}
}

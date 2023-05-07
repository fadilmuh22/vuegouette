package product

import (
	"net/http"

	"github.com/fadilmuh22/restskuy/cmd/model"
	"github.com/labstack/echo/v4"
)

// getProducts with query
func getProducts(c echo.Context) error {
	products, err := GetAllProduct()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success get all product",
		Data:    products,
	})
}

func getProduct(c echo.Context) error {
	id := c.Param("id")
	product, err := GetProduct(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success get product",
		Data:    product,
	})
}

func createProduct(c echo.Context) error {
	var product model.Product
	c.Bind(&product)

	product, err := CreateProduct(product)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success create product",
		Data:    product,
	})
}

func updateProduct(c echo.Context) error {
	var product model.Product
	c.Bind(&product)

	id := c.Param("id")
	product, err := UpdateProduct(id, product)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success update product",
		Data:    product,
	})
}

func deleteProduct(c echo.Context) error {
	id := c.Param("id")
	err := DeleteProduct(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &model.BasicResponse{
		Success: true,
		Message: "Success delete product",
		Data:    nil,
	})
}

func HandleRoutes(g *echo.Group) {
	product := g.Group("/product")
	{
		product.GET("", getProducts)
		product.POST("", createProduct)
		product.GET("/:id", getProduct)
		product.PUT("/:id", updateProduct)
		product.DELETE("/:id", deleteProduct)
	}
}

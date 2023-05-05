package product

import (
	"net/http"

	"github.com/labstack/echo"
)

// getProducts with query
func getProducts(c echo.Context) error {
	products, err := GetAllProduct()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	id := c.Param("id")
	product, err := GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func createProduct(c echo.Context) error {
	var product Product
	c.Bind(&product)

	product, err := CreateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func updateProduct(c echo.Context) error {
	var product Product
	c.Bind(&product)

	id := c.Param("id")
	product, err := UpdateProduct(id, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func deleteProduct(c echo.Context) error {
	id := c.Param("id")
	err := DeleteProduct(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Product deleted")
}

func HandleProduct(g *echo.Group) {
	g.GET("/product", getProducts)
	g.POST("/product", createProduct)
	g.GET("/product/:id", getProduct)
	g.PUT("/product/:id", updateProduct)
	g.DELETE("/product/:id", deleteProduct)
}

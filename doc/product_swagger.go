package doc

import "github.com/fadilmuh22/restskuy/cmd/model"

// swagger:route GET /product product listProducts
// List all products.
// responses:
//
//	200: productsResponse
//
// swagger:response productsResponse
type ProductsResponse struct {
	// in:body
	Body []model.Product
}

// swagger:route GET /product/{id} product getProduct
// Get a product by id.
// responses:
//
//	200: productResponse
//
// swagger:response productResponse
type ProductResponse struct {
	// in: body
	Body model.Product
}

// swagger:route POST /product product createProduct
// Create a new product.
// responses:
//
//	200: productResponse
//
// swagger:parameters createProduct
type ProductBody struct {
	// in:body
	Body model.Product
}

// swagger:route PUT /product/{id} product updateProduct
// Update a product by id.
// responses:
//
//	200: productResponse
//
// swagger:parameters updateProduct
type ProductBodyParams struct {
	// in:path
	// required:true
	ID int `json:"id"`
	// in:body
	Body model.Product
}

// swagger:route DELETE /product/{id} product deleteProduct
// Delete a product by id.
// responses:
//
//	200: productResponse
//
// swagger:parameters getProduct deleteProduct
type ProductParams struct {
	// in:path
	// required:true
	ID int `json:"id"`
}

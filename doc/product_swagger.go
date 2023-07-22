package doc

import (
	uuid "github.com/satori/go.uuid"

	"github.com/fadilmuh22/restskuy/internal/model"
)

type ProductClean struct {
	model.Product
	User struct {
		model.User
		Products []model.Product `json:"-"`
	} `json:"user"`
}

type ProductBody struct {
	model.Product
	ID     uuid.UUID  `json:"-"`
	UserID uuid.UUID  `json:"-"`
	User   model.User `json:"-"`
}

// swagger:route GET /product product listProducts
// List all products.
// responses:
//
//	200: productsResponse
//
// swagger:response productsResponse
type ProductsResponse struct {
	// in:body
	Body struct {
		model.BasicResponse
		Data []ProductClean
	} `json:"data"`
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
	Body struct {
		model.BasicResponse
		Data ProductClean `json:"data"`
	}
}

// swagger:route POST /product product createProduct
// Create a new product.
// responses:
//
//	200: productResponse
//
// swagger:parameters createProduct
type ProductCreateBody struct {
	// in:body
	Body ProductBody
}

// swagger:route PUT /product/{id} product updateProduct
// Update a product by id.
// responses:
//
//	200: productResponse
//
// swagger:parameters updateProduct
type ProductUpdateBody struct {
	// in:path
	// required:true
	ID uuid.UUID `json:"id"`
	// in:body
	Body ProductBody
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

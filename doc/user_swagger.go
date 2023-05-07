package doc

import (
	"github.com/fadilmuh22/restskuy/cmd/model"
)

// swagger:route GET /user user listUsers
// List all users.
// responses:
//
//	200: usersResponse
//
// swagger:response usersResponse
type UsersResponse struct {
	// in:body
	Body []model.User
}

// swagger:route GET /user/{id} user getUser
// Get a user by id.
// responses:
//
//	200: userResponse
//
// swagger:response userResponse
type UserResponse struct {
	// The error message
	// in: body
	Body model.User
}

// swagger:route POST /user user createUser
// Create a new user.
// responses:
//
//	200: userResponse
//
// swagger:parameters createUser
type UserBody struct {
	// in:body
	Body model.User
}

// swagger:route PUT /user/{id} user updateUser
// Update a user by id.
// responses:
//
//	200: userResponse
//
// swagger:parameters updateUser
type UserBodyParams struct {
	// in:path
	// required:true
	ID int `json:"id"`
	// in:body
	Body model.User
}

// swagger:route DELETE /user/{id} user deleteUser
// Delete a user by id.
// responses:
//
//	200: userResponse
//
// swagger:parameters getUser deleteUser
type UserParams struct {
	// in:path
	// required:true
	ID int `json:"id"`
}

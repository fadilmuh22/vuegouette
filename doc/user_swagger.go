package doc

import (
	uuid "github.com/satori/go.uuid"

	"github.com/fadilmuh22/restskuy/internal/model"
)

type UserClean struct {
	model.User
	Products []model.Product `json:"-"`
}

type UserBody struct {
	model.User
	ID       uuid.UUID       `json:"-"`
	Products []model.Product `json:"-"`
}

// swagger:route GET /user user listUsers
// List all users.
// responses:
//
//	200: usersResponse
//
// swagger:response usersResponse
type UsersResponse struct {
	// in:body
	Body struct {
		model.BasicResponse
		Data []UserClean `json:"data"`
	}
}

// swagger:route GET /user/{id} user getUser
// Get a user by id.
// responses:
//
//	200: userResponse
//
// swagger:response userResponse
type UserResponse struct {
	// in: body
	Body struct {
		model.BasicResponse
		Data UserClean `json:"data"`
	}
}

// swagger:route POST /user user createUser
// Create a new user.
// responses:
//
//	200: userResponse
//
// swagger:parameters createUser
type UserCreateBody struct {
	// in:body
	Body UserBody
}

// swagger:route PUT /user/{id} user updateUser
// Update a user by id.
// responses:
//
//	200: userResponse
//
// swagger:parameters updateUser
type UserUpdateBody struct {
	// in:path
	// required:true
	ID uuid.UUID `json:"id"`
	// in:body
	Body UserBody
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
	ID uuid.UUID `json:"id"`
}

package doc

import (
	"github.com/fadilmuh22/restskuy/internal/model"
)

// swagger:route POST /auth/login auth login
// Login.
// responses:
//
//	200: loginRegisterResponse
//
// swagger:parameters login
type LoginBody struct {
	// in:body
	Body struct {
		// The username
		// Required: true
		Email string `json:"email"`
		// The password
		// Required: true
		Password string `json:"password"`
	}
}

// swagger:route POST /auth/register auth register
// Register.
// responses:
//
//	200: loginRegisterResponse
//
// swagger:parameters register
type RegisterBody struct {
	// in:body
	Body UserBody
}

// swagger:response loginRegisterResponse
type LoginRegisterResponse struct {
	// in:body
	Body struct {
		model.BasicResponse
		Data struct {
			// The token
			Token string `json:"token"`
			// The user
			UserClean
		} `json:"data"`
	}
}

// swagger:route GET /auth/me auth me
// Get current user.
// responses:
//
//	200: meResponse
//
// swagger:response meResponse
type MeResponse struct {
	// in:body
	Body struct {
		model.BasicResponse
		Data struct {
			// The user
			UserClean
		} `json:"data"`
	}
}

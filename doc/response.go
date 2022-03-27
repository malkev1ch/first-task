package doc

import (
	"bytes"

	"github.com/malkev1ch/first-task/internal/model"
)

// swagger:parameters CreateCat
type CreateCatParam struct {
	// in:body
	// required:true
	Body model.CreateCat `json:"body"`
}

// swagger:parameters GetCat UpdateCat DeleteCat UploadCatImage
type CatUUIDParam struct {
	// in:path
	// required:true
	CatID string `json:"uuid"`
}

// swagger:response getCatResponse
type GetCatResponse struct {
	// The response message
	// in: body
	Body model.Cat `json:"body"`
}

// A GenericError is the default error message that is generated.
// For certain status codes there are more appropriate error structures.
//
// swagger:response genericError
type GenericError struct {
	// The response message
	// in: body
	Body ErrorResponseBody `json:"body"`
}

type ErrorResponseBody struct {
	// a human readable version of the error
	// required: true
	Message string `json:"message"`

	// Error An optional detailed description of the actual error. Only included if running in developer mode.
	Error string `json:"error"`
}

// swagger:parameters UploadCatImage
type UploadCatImageParam struct {
	// MyFormFile desc.
	//
	// in:formData
	//
	// swagger:file
	MyFormFile *bytes.Buffer `json:"image"`
}

// swagger:parameters SignUp
type SignUpParam struct {
	// in:body
	// required:true
	Body model.CreateUser `json:"body"`
}

// A SignUpResponse returns a couple of token with user id.
//
// swagger:response signUpResponse
type SignUpResponse struct {
	// The response message
	// in: body
	Body SignUpResponseBody `json:"body"`
}

type SignUpResponseBody struct {
	Tokens model.Tokens `json:"tokens"`
}

// swagger:parameters SignIn
type SignInParam struct {
	// in:body
	// required:true
	Body model.AuthUser `json:"body"`
}

// A SignInResponse returns a couple of token with user id.
//
// swagger:response signInResponse
type SignInResponse struct {
	// The response message
	// in: body
	Body SignInResponseBody `json:"body"`
}

type SignInResponseBody struct {
	Tokens model.Tokens `json:"tokens"`
}

// swagger:parameters RefreshToken
type RefreshTokenParam struct {
	// in:body
	// required:true
	Body model.AuthUser `json:"body"`
}

// A RefreshTokenResponse returns a couple of token with user id.
//
// swagger:response refreshTokenResponse
type RefreshTokenResponse struct {
	// The response message
	// in: body
	Body SignUpResponseBody `json:"body"`
}

type RefreshTokenResponseBody struct {
	Tokens model.Tokens `json:"tokens"`
}

// InternalServerError is a general error indicating something went wrong internally.
//
// swagger:response internalServerError
type InternalServerError GenericError

// NotFoundError is returned when the request is invalid and it cannot be processed.
//
// swagger:response notFoundError
type NotFoundError GenericError

// BadRequestError is returned when the request is invalid and it cannot be processed.
//
// swagger:response badRequestError
type BadRequestError GenericError

// An OKResponse is returned if the request was successful.
//
// swagger:response okResponse
type OKResponse struct {
	// in: body
	Body SuccessResponseBody `json:"body"`
}

type SuccessResponseBody struct {
	Message string `json:"message,omitempty"`
}

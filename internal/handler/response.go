package handler

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
	//MyFormFile desc.
	//
	// in:formData
	//
	// swagger:file
	MyFormFile *bytes.Buffer `json:"image"`
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

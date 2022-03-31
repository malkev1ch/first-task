package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

//	swagger:route POST /auth/sign-up auth SignUp
//
//	Registration process for new user
//
//	Returns a couple of tokens for recently created user.
//
//	responses:
//	 201: signUpResponse
//	 400: badRequestError
//	 415: unsupportedMediaTypeError
//	 500: internalServerError
func (h *Handler) SignUp(ctx echo.Context) error {
	contentType := ctx.Request().Header.Get("Content-Type")
	if _, ex := AllowedContentType[contentType]; !ex {
		return ctx.JSON(http.StatusUnsupportedMediaType, ErrorResponse{
			Message: "invalid media type", Error: fmt.Sprintf("got - %s", contentType),
		})
	}

	var input model.CreateUser
	if err := ctx.Bind(&input); err != nil {
		logrus.Error("handler: invalid content of body - ", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid body content", Error: err.Error(),
		})
	}

	if err := h.Validator.Validate(&input); err != nil {
		logrus.Error("handler: validation failed - ", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "not enough fields in json body or wrong values of fields", Error: err.Error(),
		})
	}

	tokens, err := h.Services.SignUp(ctx.Request().Context(), &input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't create user", Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, *tokens)
}

//	swagger:route POST /auth/sign-in auth SignIn
//
//	Authorisation process for existed user
//
//	Returns a couple of tokens for existed user.
//
//	responses:
//	 200: signInResponse
//	 400: badRequestError
//	 415: unsupportedMediaTypeError
//	 500: internalServerError
func (h *Handler) SignIn(ctx echo.Context) error {
	contentType := ctx.Request().Header.Get("Content-Type")
	if _, ex := AllowedContentType[contentType]; !ex {
		return ctx.JSON(http.StatusUnsupportedMediaType, ErrorResponse{
			Message: "invalid media type", Error: fmt.Sprintf("got - %s", contentType),
		})
	}

	var input model.AuthUser
	if err := ctx.Bind(&input); err != nil {
		logrus.Error("handler: invalid content of body - ", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid body content", Error: err.Error(),
		})
	}

	if err := h.Validator.Validate(&input); err != nil {
		logrus.Error("handler: validation failed - ", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "not enough fields in json body or wrong values of fields", Error: err.Error(),
		})
	}

	tokens, err := h.Services.SignIn(ctx.Request().Context(), &input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "authorisation failed", Error: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, *tokens)
}

//	swagger:route POST /auth/refresh auth RefreshToken
//
//	Refresh a couple tokens used refresh token
//
//	Returns a couple tokens for existed user.
//
//	responses:
//	 200: refreshTokenResponse
//	 400: badRequestError
//	 415: unsupportedMediaTypeError
//	 500: internalServerError
func (h *Handler) RefreshToken(ctx echo.Context) error {
	contentType := ctx.Request().Header.Get("Content-Type")
	if _, ex := AllowedContentType[contentType]; !ex {
		return ctx.JSON(http.StatusUnsupportedMediaType, ErrorResponse{
			Message: "invalid media type", Error: fmt.Sprintf("got - %s", contentType),
		})
	}

	var input model.RefreshToken
	if err := ctx.Bind(&input); err != nil {
		logrus.Error("handler: invalid content of body - ", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid body content", Error: err.Error(),
		})
	}

	if err := h.Validator.Validate(&input); err != nil {
		logrus.Error("handler: validation failed - ", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "not enough fields", Error: err.Error(),
		})
	}

	tokens, err := h.Services.RefreshToken(ctx.Request().Context(), input.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "failed refresh token", Error: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, *tokens)
}

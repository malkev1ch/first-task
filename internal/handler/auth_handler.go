package handler

import (
	"net/http"

	"github.com/malkev1ch/first-task/internal/service"

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
//	 500: internalServerError
func (h *Handler) SignUp(ctx echo.Context) error {
	var input model.CreateUser
	if err := ctx.Bind(&input); err != nil {
		logrus.Error(err, "handler: invalid content of body")
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid body content", Error: err.Error(),
		})
	}

	tokens, err := h.Services.SignUp(ctx.Request().Context(), &service.SignUpInput{
		Email: input.Email, Password: input.Password,
		UserName: input.UserName,
	})
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
//	 500: internalServerError
func (h *Handler) SignIn(ctx echo.Context) error {
	var input model.AuthUser
	if err := ctx.Bind(&input); err != nil {
		logrus.Error(err, "handler: invalid content of body")
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid body content", Error: err.Error(),
		})
	}

	tokens, err := h.Services.SignIn(ctx.Request().Context(), &service.SignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
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
//	 500: internalServerError
func (h *Handler) RefreshToken(ctx echo.Context) error {
	var input model.RefreshToken
	if err := ctx.Bind(&input); err != nil {
		logrus.Error(err, "handler: invalid content of body")
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid body content", Error: err.Error(),
		})
	}

	tokens, err := h.Services.RefreshToken(ctx.Request().Context(), input.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "authorisation failed", Error: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, *tokens)
}
package handler

import (
	"github.com/labstack/echo"
	"github.com/malkev1ch/first-task/internal/domain"
	"net/http"
	"strconv"
)

//----------
// Handlers
//----------

func (h *Handler) CreateCat(ctx echo.Context) error {
	var input domain.CreateCat
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.services.Cat.CreateCat(input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, dataResponse{
		Data: id,
	})
}

func (h *Handler) UpdateCat(ctx echo.Context) error {

	var input domain.UpdateCat

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.services.Cat.UpdateCat(id, input); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusInternalServerError, err.Error())
}

func (h *Handler) DeleteCat(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.services.Cat.DeleteCat(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "The dish has been successfully deleted")
}

func (h *Handler) GetCat(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	cat, err := h.services.Cat.GetCat(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dataResponse{
		Data: cat,
	})
}

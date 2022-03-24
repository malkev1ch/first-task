package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

func (h *Handler) CreateCat(ctx echo.Context) error {
	var input model.Cat
	if err := ctx.Bind(&input); err != nil {
		logrus.Error(fmt.Errorf("handler: can't create cat - %w", err))
		return ctx.String(http.StatusInternalServerError, "error while creating cat")
	}

	id, err := h.Services.Create(ctx.Request().Context(), &input)
	if err != nil {
		logrus.Error(fmt.Errorf("handler: can't create cat - %w", err))
		return ctx.String(http.StatusInternalServerError, "error while creating cat")
	}

	type idResponse struct {
		ID string `json:"id"`
	}

	return ctx.JSON(http.StatusCreated, idResponse{
		ID: id,
	})
}

func (h *Handler) GetCat(ctx echo.Context) error {
	id := ctx.Param("id")
	cat, err := h.Services.Get(ctx.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, cat)
}

func (h *Handler) UpdateCat(ctx echo.Context) error {
	id := ctx.Param("id")
	var input model.Cat
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Services.Update(ctx.Request().Context(), id, &input); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, "")
}

func (h *Handler) DeleteCat(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := h.Services.Delete(ctx.Request().Context(), id); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "The cat has been successfully deleted")
}

func (h *Handler) UploadCatImage(ctx echo.Context) error {
	id := ctx.Param("id")
	file, err := ctx.FormFile("image")
	if err != nil {
		logrus.Errorf("handler: can't parse form file - %e", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		logrus.Errorf("handler: can't open file - %e", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	filename := h.generateFileName(file.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		logrus.Errorf("handler: can't create file locally - %e", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	buffer := make([]byte, file.Size)

	if _, err = src.Read(buffer); err != nil {
		logrus.Errorf("handler: can't read file - %e", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	contentType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := imageTypes[contentType]; !ex {
		logrus.Errorf("cannot read the file")
		return ctx.String(http.StatusBadRequest, "file type is not supported")
	}

	if _, err = io.Copy(dst, src); err != nil {
		logrus.Errorf("handler: can't copy file - %e", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	if err := h.Services.UploadImage(ctx.Request().Context(), id, filename); err != nil {
		logrus.Errorf("handler: can't update cats image path - %e", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	type response struct {
		Path string `json:"path"`
	}

	return ctx.JSON(http.StatusOK, response{
		Path: filename,
	})
}

func (h *Handler) generateFileName(filename string) string {
	resultedFilename := fmt.Sprintf("%s%s.%s", h.Cfg.ImagePath, uuid.New().String(), getFileExtension(filename))
	return resultedFilename
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}

func (h *Handler) GetCatImage(ctx echo.Context) error {
	id := ctx.Param("id")
	cat, err := h.Services.Get(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.File(*cat.ImagePath)
}

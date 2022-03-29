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
	"image/webp": nil,
}

//	swagger:route POST /cats cats CreateCat
//
//	Create cat
//
//	Creates a new cat.
//
//	Security:
//	 AdminAuth:
//
//	responses:
//	 201: okResponse
//	 400: badRequestError
//	 401: unauthorizedError
//	 500: internalServerError
func (h *Handler) CreateCat(ctx echo.Context) error {
	var input model.CreateCat
	if err := ctx.Bind(&input); err != nil {
		logrus.Error(fmt.Errorf("handler: invalid content pf body - %w", err))
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid path parameter", Error: err.Error(),
		})
	}

	id, err := h.Services.Create(ctx.Request().Context(), &model.Cat{
		Name:       input.Name,
		DateBirth:  input.DateBirth,
		Vaccinated: input.Vaccinated,
	})
	if err != nil {
		logrus.Error(fmt.Errorf("handler: can't create cat - %w", err))
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't create cat", Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, OKResponse{
		Message: id,
	})
}

//	swagger:route GET /cats/{uuid} cats GetCat
//
//	Get cat by UUID.
//
//	Returns a cat with the given UUID.
//
//	Security:
//	 AdminAuth:
//
//	responses:
//	 200: getCatResponse
//	 401: unauthorizedError
//	 500: internalServerError
func (h *Handler) GetCat(ctx echo.Context) error {
	id := ctx.Param("uuid")
	cat, err := h.Services.Get(ctx.Request().Context(), id)
	if err != nil {
		logrus.Error(err, "handler: can't get cat")
		return echo.NewHTTPError(http.StatusInternalServerError, ErrorResponse{
			Message: "can't get cat", Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, *cat)
}

//	swagger:route PUT /cats/{uuid} cats UpdateCat
//
//	Update cat.
//
//	Update a cat with the given UUID.
//
//	Security:
//	 AdminAuth:
//
//	responses:
//	 200: updateCatResponse
//	 400: badRequestError
//	 401: unauthorizedError
//	 500: internalServerError
func (h *Handler) UpdateCat(ctx echo.Context) error {
	id := ctx.Param("uuid")
	var input model.UpdateCat
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid path parameter", Error: err.Error(),
		})
	}

	cat, err := h.Services.Update(ctx.Request().Context(), id, &input)
	if err != nil {
		logrus.Error(fmt.Errorf("handler: can't create cat - %w", err))
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't update cat", Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, *cat)
}

//	swagger:route DELETE /cats/{uuid} cats DeleteCat
//
//	Remove cat from storage.
//
//	Security:
//	 AdminAuth:
//
//	Responses:
//	 200: okResponse
//	 401: unauthorizedError
//	 500: internalServerError
func (h *Handler) DeleteCat(ctx echo.Context) error {
	id := ctx.Param("uuid")
	if err := h.Services.Delete(ctx.Request().Context(), id); err != nil {
		logrus.Error(fmt.Errorf("handler: can't delete cat - %w", err))
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't delete cat", Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, OKResponse{
		Message: "OK",
	})
}

//	swagger:route POST /cats/{uuid}/image cats UploadCatImage
//
// 	Set or update cats image.
//
// 	consumes:
//   - multipart/form-data
//
//	Security:
//	 AdminAuth:
//
// 	Responses:
// 	 200: okResponse
// 	 400: badRequestError
//	 401: unauthorizedError
// 	 500: internalServerError
func (h *Handler) UploadCatImage(ctx echo.Context) error {
	id := ctx.Param("uuid")
	file, err := ctx.FormFile("image")
	if err != nil {
		logrus.Errorf("handler: can't parse form file - %e", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "can't parse form file", Error: err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		logrus.Errorf("handler: can't open file - %e", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't open file", Error: err.Error(),
		})
	}
	defer src.Close()

	filename := h.generateFileName(file.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		logrus.Errorf("handler: can't create file locally - %e", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't create file locally", Error: err.Error(),
		})
	}
	defer dst.Close()

	buffer := make([]byte, file.Size)

	if _, err = src.Read(buffer); err != nil {
		logrus.Errorf("handler: can't read file - %e", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't read file", Error: err.Error(),
		})
	}

	contentType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := imageTypes[contentType]; !ex {
		logrus.Errorf("invalid file type")
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "image should be .jpg or .png", Error: "invalid file type",
		})
	}

	if _, err = io.Copy(dst, src); err != nil {
		logrus.Errorf("handler: can't copy file - %e", err)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "can't copy file", Error: err.Error(),
		})
	}

	if err := h.Services.UploadImage(ctx.Request().Context(), id, filename); err != nil {
		logrus.Errorf("handler: can't update cats image path - %e", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't update cats image path", Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, OKResponse{
		Message: filename,
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

//	swagger:route GET /cats/{uuid}/image cats GetCatImage
//
//	Get cats image.
//
//	Produces:
//	 - image/jpeg
//	 - image/png
//	 - image/webp
//
//	Security:
//	 AdminAuth:
//
// 	Responses:
// 	 200: okResponse
// 	 500: internalServerError
func (h *Handler) GetCatImage(ctx echo.Context) error {
	id := ctx.Param("uuid")
	cat, err := h.Services.Get(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "can't update cats image path", Error: err.Error(),
		})
	}

	return ctx.File(cat.ImagePath)
}

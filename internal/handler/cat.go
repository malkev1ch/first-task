package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/malkev1ch/first-task/internal/domain"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//var (
//	imageTypes = map[string]interface{}{
//		"image/jpeg": nil,
//		"image/png":  nil,
//	}
//)

func (h *Handler) CreateCat(ctx echo.Context) error {
	var input domain.CreateCat
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.services.Cat.CreateCat(input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dataResponse{
		Data: id,
	})
}

func (h *Handler) UpdateCat(ctx echo.Context) error {

	var input domain.UpdateCat

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.services.Cat.UpdateCat(id, input); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, "")
}

func (h *Handler) DeleteCat(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.services.Cat.DeleteCat(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "The cat has been successfully deleted")
}

func (h *Handler) GetCat(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	cat, err := h.services.Cat.GetCat(id)
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dataResponse{
		Data: cat,
	})
}

func (h *Handler) uploadCatImage(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	filename := h.generateFileName(file.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		logrus.Error(err)
		fmt.Println(112)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := h.services.UploadImage(id, filename); err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dataResponse{
		Data: filename,
	})

	//TODO
	//contentType := http.DetectContentType(buffer)
	//
	//// Validate File Type
	//if _, ex := imageTypes[contentType]; !ex {
	//	logrus.Info("cannot read the file")
	//	return ctx.JSON(http.StatusBadRequest, "file type is not supported")
	//}
}

func (h * Handler) generateFileName(filename string) string {
	resultedFilename := fmt.Sprintf("%s%s.%s", h.cfg.Image.PhotoPathCat, uuid.New().String(), getFileExtension(filename))
	return resultedFilename
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}

func (h *Handler) getCatImage(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	cat, err := h.services.Cat.GetCat(id)
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.File(cat.ImagePath)
}

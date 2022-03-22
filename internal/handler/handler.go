package handler

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/malkev1ch/first-task/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}


func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.Logger.SetLevel(log.DEBUG)
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE},
	}))

	cat := router.Group("/cat")
	{
		cat.GET("/:id", h.GetCat)
		cat.POST("/", h.CreateCat)
		cat.PUT("/:id", h.UpdateCat)
		cat.DELETE("/:id", h.DeleteCat)
	}

	return router
}

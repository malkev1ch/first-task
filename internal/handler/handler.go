package handler

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/malkev1ch/first-task/internal/config"
	"github.com/malkev1ch/first-task/internal/service"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type OKResponse struct {
	Message string `json:"message"`
}

var AllowedContentType = map[string]interface{}{
	"application/json": nil,
}

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
	"image/webp": nil,
}

// Handler type replies for handling echo server requests.
type Handler struct {
	Services  *service.Service
	Cfg       *config.Config
	Validator *Validator
}

// NewHandler function create handler.
func NewHandler(services *service.Service, cfg *config.Config, validator *Validator) *Handler {
	return &Handler{
		Services:  services,
		Cfg:       cfg,
		Validator: validator,
	}
}

func InitRouter(handlers *Handler, cfg *config.Config) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", handlers.SignUp)
		auth.POST("/sign-in", handlers.SignIn)
		auth.POST("/refresh", handlers.RefreshToken)
	}

	cat := router.Group("/cats")

	if cfg.AuthMode {
		configJWTMiddleware := middleware.JWTConfig{
			Claims:     &service.JwtCustomClaims{},
			SigningKey: []byte(cfg.JWTKey),
		}
		cat.Use(middleware.JWTWithConfig(configJWTMiddleware))
	}

	{
		cat.GET("/:uuid", handlers.GetCat)
		cat.POST("/", handlers.CreateCat)
		cat.PUT("/:uuid", handlers.UpdateCat)
		cat.DELETE("/:uuid", handlers.DeleteCat)
		cat.POST("/:uuid/image", handlers.UploadCatImage)
		cat.GET("/:uuid/image", handlers.GetCatImage)
	}
	return router
}

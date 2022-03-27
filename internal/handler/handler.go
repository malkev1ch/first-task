package handler

import (
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

// Handler type replies for handling echo server requests.
type Handler struct {
	Services *service.Service
	Cfg      *config.Config
}

// NewHandler function create handler.
func NewHandler(services *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		Services: services,
		Cfg:      cfg,
	}
}

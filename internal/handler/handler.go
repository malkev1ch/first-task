package handler

import (
	"github.com/malkev1ch/first-task/internal/configs"
	"github.com/malkev1ch/first-task/internal/service"
)

type Handler struct {
	Services *service.Service
	Cfg      *configs.Config
}

func NewHandler(services *service.Service, cfg *configs.Config) *Handler {
	return &Handler{Services: services,
		Cfg: cfg}
}

package server

import (
	"context"
	"github.com/malkev1ch/first-task/configs"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg *configs.Config, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

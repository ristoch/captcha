package http

import (
	"context"
	"net/http"
	"time"
)

type BalancerServer struct {
	handlers *BalancerHandlers
	port     string
	server   *http.Server
}

func NewBalancerServer(handlers *BalancerHandlers, port string) *BalancerServer {
	return &BalancerServer{
		handlers: handlers,
		port:     port,
	}
}

func (s *BalancerServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handlers.HealthHandler)
	mux.HandleFunc("/api/health", s.handlers.APIHealthHandler)
	mux.HandleFunc("/api/services", s.handlers.ServicesHandler)

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	return s.server.ListenAndServe()
}

func (s *BalancerServer) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

package http

import (
	"context"
	"net/http"

	"captcha-service/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	handlers *Handlers
	port     string
}

func NewServer(handlers *Handlers, port string) *Server {
	return &Server{
		handlers: handlers,
		port:     port,
	}
}

func (s *Server) Start() error {
	r := mux.NewRouter()

	r.HandleFunc("/api/challenge", s.handlers.HandleChallengeRequest).Methods("POST")
	r.HandleFunc("/ws", s.handlers.HandleWebSocket)
	r.HandleFunc("/health", s.handlers.HandleHealthCheck).Methods("GET")

	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: r,
	}

	logger.Info("HTTP server starting", zap.String("port", s.port))
	return server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	// HTTP server doesn't have a graceful stop method in this implementation
	// In a real implementation, you would use server.Shutdown(ctx)
	return nil
}

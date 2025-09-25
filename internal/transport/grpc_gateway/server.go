package grpc_gateway

import (
	"context"
	"fmt"
	"net/http"

	grpcTransport "captcha-service/internal/transport/grpc"
	httpTransport "captcha-service/internal/transport/http"
	"captcha-service/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	grpcHandlers *grpcTransport.Handlers
	httpHandlers *httpTransport.Handlers
	port         int
	httpServer   *http.Server
}

func NewServer(grpcHandlers *grpcTransport.Handlers, httpHandlers *httpTransport.Handlers, port int) *Server {
	return &Server{
		grpcHandlers: grpcHandlers,
		httpHandlers: httpHandlers,
		port:         port,
	}
}

func (s *Server) Start() error {
	// Создаем HTTP роутер
	router := mux.NewRouter()

	// Добавляем HTTP маршруты напрямую
	router.HandleFunc("/api/challenge", s.httpHandlers.HandleChallengeRequest).Methods("POST")
	router.HandleFunc("/api/validate", s.httpHandlers.HandleValidateRequest).Methods("POST")
	router.HandleFunc("/ws", s.httpHandlers.HandleWebSocket)
	router.HandleFunc("/health", s.httpHandlers.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/memory", s.httpHandlers.HandleMemoryStats).Methods("GET")
	router.HandleFunc("/stats", s.httpHandlers.HandleStats).Methods("GET")

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router,
	}

	logger.Info("gRPC-Gateway server starting", zap.Int("port", s.port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

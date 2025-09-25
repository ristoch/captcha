package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	captchav1 "captcha-service/gen/proto/captcha"

	"google.golang.org/grpc"
)

type Server struct {
	handlers *Handlers
	port     int
	server   *grpc.Server
}

func NewServer(handlers *Handlers, port int) *Server {
	return &Server{
		handlers: handlers,
		port:     port,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	captchav1.RegisterCaptchaServiceServer(s.server, s.handlers)

	log.Printf("gRPC server starting on port %d", s.port)
	return s.server.Serve(lis)
}

func (s *Server) Register(grpcServer *grpc.Server) {
	captchav1.RegisterCaptchaServiceServer(grpcServer, s.handlers)
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

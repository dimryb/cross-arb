package grpc

import (
	"fmt"
	"net"

	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	app i.Application
	cfg ServerConfig
	log i.Logger
}

type ServerConfig struct {
	Port string
}

func NewServer(app i.Application, cfg ServerConfig, log i.Logger) *Server {
	return &Server{
		app: app,
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// TODO: добавить интерцепторы (логирование, метрики и т.п.)
	grpcServer := grpc.NewServer(
	// grpc.UnaryInterceptor(...),
	)

	// Регистрируем сервис
	proto.RegisterTickerServiceServer(grpcServer, NewTickerService(s.app))

	// Включаем reflection — удобно для CLI (grpcurl, evans)
	reflection.Register(grpcServer)

	s.log.Infof("Starting gRPC server on port %s", s.cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

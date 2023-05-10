package internalgrpc

import (
	"context"
	"log"
	"net"

	router "github.com/BBeRsErKeRR/system-stats-monitor/api"
	v1grpc "github.com/BBeRsErKeRR/system-stats-monitor/api/v1/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	pkgnet "github.com/BBeRsErKeRR/system-stats-monitor/pkg/net"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	logger logger.Logger
	Addr   string
	server *grpc.Server
}

type Config struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewServer(logger logger.Logger, app router.Application, conf *Config) *Server {
	addr, err := pkgnet.GetAddress(conf.Host, conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingMiddleware(logger),
		),
	)
	v1grpc.RegisterSystemStatsMonitorServiceV1Server(server, v1grpc.NewHandler(app, logger))
	return &Server{
		Addr:   addr,
		logger: logger,
		server: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	list, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.logger.Info("starting server", zap.String("address", s.Addr))
	err = s.server.Serve(list)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop() error {
	s.server.GracefulStop()
	return nil
}

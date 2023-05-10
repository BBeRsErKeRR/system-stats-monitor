package internalgrpc

import (
	"context"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func loggingMiddleware(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		r interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		response, err := handler(ctx, r)
		status := "success"
		if err != nil {
			status = "error"
		}
		logger.Info(
			"grpc-access-log",
			zap.String("datetime", time.Now().Format(time.RFC822)),
			zap.String("method", info.FullMethod),
			zap.Any("duration", time.Since(start).String()),
			zap.String("status", status),
		)
		return response, err
	}
}

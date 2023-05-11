package v1grpc

import (
	router "github.com/BBeRsErKeRR/system-stats-monitor/api"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"go.uber.org/zap"
)

type Handler struct {
	app    router.Application
	logger logger.Logger
	UnimplementedSystemStatsMonitorServiceV1Server
}

func NewHandler(app router.Application, logger logger.Logger) *Handler {
	return &Handler{
		app:    app,
		logger: logger,
	}
}

func getStatsResponse(stats stats.Stats) *StatsResponse {
	cpuInfo := &CPUTimeStatValue{
		User:   stats.CPUInfo.User,
		System: stats.CPUInfo.System,
		Idle:   stats.CPUInfo.Idle,
	}

	return &StatsResponse{
		CpuInfo: cpuInfo,
	}
}

func (h *Handler) StartMonitoring(req *StartMonitoringRequest, srv SystemStatsMonitorServiceV1_StartMonitoringServer) error { //nolint:lll
	result, err := h.app.StartMonitoring(srv.Context(), req.GetResponseDuration(), req.GetWaitDuration())
	if err != nil {
		return err
	}
	for stats := range result {
		resp := getStatsResponse(stats)
		if err := srv.Send(resp); err != nil {
			h.logger.Error("send error", zap.Error(err))
		}
	}
	return nil
}

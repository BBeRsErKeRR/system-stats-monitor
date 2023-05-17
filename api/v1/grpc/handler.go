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

	loadInfo := &LoadStatValue{
		Load1:  stats.LoadInfo.Load1,
		Load5:  stats.LoadInfo.Load5,
		Load15: stats.LoadInfo.Load15,
	}

	nsInfo := &NetworkStateStatValue{
		Counters: stats.NetworkStateInfo.Counters,
	}

	statisticsItems := make([]*NetworkStatItem, 0, len(stats.NetworkStatisticsInfo.Items))
	for _, statItem := range stats.NetworkStatisticsInfo.Items {
		statisticsItems = append(statisticsItems, &NetworkStatItem{
			Command:  statItem.Command,
			Pid:      statItem.PID,
			User:     statItem.User,
			Protocol: statItem.Protocol,
			Port:     statItem.Port,
		})
	}

	duItems := make([]*DiskUsageItem, 0, len(stats.DiskUsageInfo.Items))
	for _, statItem := range stats.DiskUsageInfo.Items {
		duItems = append(duItems, &DiskUsageItem{
			Path:                   statItem.Path,
			Fstype:                 statItem.Fstype,
			Used:                   statItem.Used,
			AvailablePercent:       statItem.AvailablePercent,
			InodeUsed:              statItem.InodesUsed,
			InodesAvailablePercent: statItem.InodesAvailablePercent,
		})
	}

	dIoItems := make([]*DiskIoItem, 0, len(stats.DiskIoInfo.Items))
	for _, statItem := range stats.DiskIoInfo.Items {
		dIoItems = append(dIoItems, &DiskIoItem{
			Device:   statItem.Device,
			Tps:      statItem.Tps,
			KbReadS:  statItem.KbReadS,
			KbWriteS: statItem.KbWriteS,
		})
	}

	protocolTalkers := make([]*ProtocolTalkerItem, 0, len(stats.ProtocolTalkersInfo.Items))
	for _, statItem := range stats.ProtocolTalkersInfo.Items {
		protocolTalkers = append(protocolTalkers, &ProtocolTalkerItem{
			Protocol:        statItem.Protocol,
			SendBytes:       statItem.SendBytes,
			BytesPercentage: statItem.BytesPercentage,
		})
	}
	bpsTalkers := make([]*BpsTalkerItem, 0, len(stats.BpsTalkersInfo.Items))
	for _, statItem := range stats.BpsTalkersInfo.Items {
		bpsTalkers = append(bpsTalkers, &BpsTalkerItem{
			Source:      statItem.Source,
			Destination: statItem.Destination,
			Protocol:    statItem.Protocol,
			Bps:         statItem.Bps,
			Numbers:     statItem.Numbers,
		})
	}

	return &StatsResponse{
		CpuInfo:          cpuInfo,
		LoadInfo:         loadInfo,
		NetworkStateInfo: nsInfo,
		NetworkStatisticsInfo: &NetworkStatisticsValue{
			Items: statisticsItems,
		},
		DiskUsageInfo: &DiskUsageValue{
			Items: duItems,
		},
		DiskIoInfo: &DiskIoValue{
			Items: dIoItems,
		},
		ProtocolTalkers: &ProtocolTalkersValue{
			Items: protocolTalkers,
		},
		BpsTalkers: &BpsTalkersValue{
			Items: bpsTalkers,
		},
	}
}

func (h *Handler) StartMonitoring(req *StartMonitoringRequest, srv SystemStatsMonitorServiceV1_StartMonitoringServer) error { //nolint:lll
	st := h.app.CreateStorage()
	ctx := srv.Context()
	err := st.Connect(ctx)
	if err != nil {
		return err
	}
	defer st.Close(ctx)
	useCase := h.app.CreateUseCase(st)
	result, err := h.app.StartMonitoring(ctx, req.GetResponseDuration(), req.GetWaitDuration(), useCase)
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

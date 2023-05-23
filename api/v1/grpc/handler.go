package v1grpc

import (
	router "github.com/BBeRsErKeRR/system-stats-monitor/api"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
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

func ResolveResponse(resp *StatsResponse) *stats.Stats {
	statisticsItems := make([]storage.NetworkStatsItem, 0, len(resp.GetNetworkStatisticsInfo().GetItems()))
	for _, statItem := range resp.GetNetworkStatisticsInfo().GetItems() {
		statisticsItems = append(statisticsItems, storage.NetworkStatsItem{
			Command:  statItem.Command,
			PID:      statItem.Pid,
			User:     statItem.User,
			Protocol: statItem.Protocol,
			Port:     statItem.Port,
		})
	}

	duItems := make([]storage.UsageStatItem, 0, len(resp.GetDiskUsageInfo().GetItems()))
	for _, respItem := range resp.GetDiskUsageInfo().GetItems() {
		duItems = append(duItems, storage.UsageStatItem{
			Path:                   respItem.Path,
			Fstype:                 respItem.Fstype,
			Used:                   respItem.Used,
			AvailablePercent:       respItem.AvailablePercent,
			InodesUsed:             respItem.InodeUsed,
			InodesAvailablePercent: respItem.InodesAvailablePercent,
		})
	}

	dIoItems := make([]storage.DiskIoStatItem, 0, len(resp.GetDiskIoInfo().GetItems()))
	for _, statItem := range resp.GetDiskIoInfo().GetItems() {
		dIoItems = append(dIoItems, storage.DiskIoStatItem{
			Device:   statItem.Device,
			Tps:      statItem.Tps,
			KbReadS:  statItem.KbReadS,
			KbWriteS: statItem.KbWriteS,
		})
	}

	protocolTalkers := make([]storage.ProtocolTalkerItem, 0, len(resp.GetProtocolTalkers().GetItems()))
	for _, statItem := range resp.GetProtocolTalkers().GetItems() {
		protocolTalkers = append(protocolTalkers, storage.ProtocolTalkerItem{
			Protocol:        statItem.Protocol,
			SendBytes:       statItem.SendBytes,
			BytesPercentage: statItem.BytesPercentage,
		})
	}
	bpsTalkers := make([]storage.BpsItem, 0, len(resp.GetBpsTalkers().GetItems()))
	for _, statItem := range resp.GetBpsTalkers().GetItems() {
		bpsTalkers = append(bpsTalkers, storage.BpsItem{
			Source:      statItem.Source,
			Destination: statItem.Destination,
			Protocol:    statItem.Protocol,
			Bps:         statItem.Bps,
			Numbers:     statItem.Numbers,
		})
	}
	return &stats.Stats{
		CPUInfo: storage.CPUTimeStat{
			User:   resp.GetCpuInfo().GetUser(),
			System: resp.GetCpuInfo().GetSystem(),
			Idle:   resp.GetCpuInfo().GetIdle(),
		},
		LoadInfo: storage.LoadStat{
			Load1:  resp.GetLoadInfo().GetLoad1(),
			Load5:  resp.GetLoadInfo().GetLoad5(),
			Load15: resp.GetLoadInfo().GetLoad15(),
		},
		NetworkStateInfo: storage.NetworkStatesStat{
			Counters: resp.GetNetworkStateInfo().GetCounters(),
		},
		NetworkStatisticsInfo: storage.NetworkStats{
			Items: statisticsItems,
		},
		DiskUsageInfo: storage.UsageStats{
			Items: duItems,
		},
		DiskIoInfo: storage.DiskIoStat{
			Items: dIoItems,
		},
		ProtocolTalkersInfo: storage.ProtocolTalkersStats{
			Items: protocolTalkers,
		},
		BpsTalkersInfo: storage.BpsTalkersStats{
			Items: bpsTalkers,
		},
	}
}

func (h *Handler) StartMonitoring(req *StartMonitoringRequest, srv SystemStatsMonitorServiceV1_StartMonitoringServer) error { //nolint:lll
	st := h.app.CreateStorage()
	ctx := srv.Context()
	if err := ctx.Err(); err != nil {
		return err
	}
	err := st.Connect(ctx)
	if err != nil {
		return err
	}
	defer st.Close(ctx)
	useCase, err := h.app.CreateUseCase(ctx, st)
	if err != nil {
		return err
	}
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

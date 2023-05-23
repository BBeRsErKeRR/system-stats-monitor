package stats

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/cpu"
	diskio "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/disk/io"
	diskusage "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/disk/usage"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/load"
	networkstates "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/network/states"
	networkstatistics "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/network/statistics"
	bpstalkers "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/network/talkers/bps"
	protocoltalkers "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/network/talkers/protocol"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"go.uber.org/zap"
)

var (
	ErrCollector = errors.New("unsupported collector type")
	ErrStatsType = errors.New("unsupported stats type")
)

type Config struct {
	IsCPUEnable            bool `mapstructure:"cpu_enable"`
	IsLoadEnable           bool `mapstructure:"load_enable"`
	IsNetworkEnable        bool `mapstructure:"network_enable"`
	IsDiskEnable           bool `mapstructure:"disk_enable"`
	IsNetworkTalkersEnable bool `mapstructure:"network_talkers_enable"`
}

type Stats struct {
	CPUInfo               storage.CPUTimeStat          `json:"cpu_info"`                //nolint:tagliatelle
	LoadInfo              storage.LoadStat             `json:"load_info"`               //nolint:tagliatelle
	NetworkStateInfo      storage.NetworkStatesStat    `json:"network_state_info"`      //nolint:tagliatelle
	NetworkStatisticsInfo []storage.NetworkStatsItem   `json:"network_statistics_info"` //nolint:tagliatelle
	DiskUsageInfo         []storage.UsageStatItem      `json:"disk_usage_info"`         //nolint:tagliatelle
	DiskIoInfo            []storage.DiskIoStatItem     `json:"disk_io_info"`            //nolint:tagliatelle
	ProtocolTalkersInfo   []storage.ProtocolTalkerItem `json:"protocol_talkers"`        //nolint:tagliatelle
	BpsTalkersInfo        []storage.BpsItem            `json:"bps_talkers"`             //nolint:tagliatelle
}

type UseCase struct {
	logger                 logger.Logger
	st                     storage.Storage
	collectors             map[string]monitor.Collector
	constantCollectors     map[string]monitor.ConstantCollector
	isCPUEnable            bool
	isLoadEnable           bool
	isNetworkEnable        bool
	isDiskEnable           bool
	isNetworkTalkersEnable bool
}

func CheckExecution(ctx context.Context, cfg *Config, logger logger.Logger) error {
	if cfg.IsCPUEnable {
		if err := cpu.New(logger).CheckExecution(ctx); err != nil {
			return err
		}
	}

	if cfg.IsLoadEnable {
		if err := load.New(logger).CheckExecution(ctx); err != nil {
			return err
		}
	}

	if cfg.IsNetworkEnable {
		if err := networkstates.New(logger).CheckExecution(ctx); err != nil {
			return err
		}

		if err := networkstatistics.New(logger).CheckExecution(ctx); err != nil {
			return err
		}
	}

	if cfg.IsDiskEnable {
		if err := diskusage.New(logger).CheckExecution(ctx); err != nil {
			return err
		}

		if err := diskio.New(logger).CheckExecution(ctx); err != nil {
			return err
		}
	}

	if cfg.IsNetworkTalkersEnable {
		if err := protocoltalkers.New(logger).CheckExecution(ctx); err != nil {
			return err
		}
		if err := bpstalkers.New(logger).CheckExecution(ctx); err != nil {
			return err
		}
	}

	return nil
}

func New(ctx context.Context, cfg *Config, st storage.Storage, logger logger.Logger) (UseCase, error) {
	collectors := map[string]monitor.Collector{}
	constantCollectors := map[string]monitor.ConstantCollector{}

	if cfg.IsCPUEnable {
		cpuC := cpu.New(logger)
		if err := cpuC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		collectors["cpu"] = cpuC
	}

	if cfg.IsLoadEnable {
		loadC := load.New(logger)
		if err := loadC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		collectors["load"] = loadC
	}

	if cfg.IsNetworkEnable {
		nstC := networkstates.New(logger)
		if err := nstC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		collectors["network_states"] = nstC

		nsC := networkstatistics.New(logger)
		if err := nsC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		collectors["network_statistics"] = nsC
	}

	if cfg.IsDiskEnable {
		duC := diskusage.New(logger)
		if err := duC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		collectors["du"] = duC

		dioC := diskio.New(logger)
		if err := dioC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		collectors["di"] = dioC
	}

	if cfg.IsNetworkTalkersEnable {
		protocolC := protocoltalkers.New(logger)
		if err := protocolC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		constantCollectors["protocol_talkers"] = protocolC
		bspC := bpstalkers.New(logger)
		if err := bspC.CheckExecution(ctx); err != nil {
			return UseCase{}, err
		}
		constantCollectors["bsp_talkers"] = bspC
	}

	return UseCase{
		collectors:             collectors,
		constantCollectors:     constantCollectors,
		st:                     st,
		logger:                 logger,
		isCPUEnable:            cfg.IsCPUEnable,
		isLoadEnable:           cfg.IsLoadEnable,
		isNetworkEnable:        cfg.IsNetworkEnable,
		isDiskEnable:           cfg.IsDiskEnable,
		isNetworkTalkersEnable: cfg.IsNetworkTalkersEnable,
	}, nil
}

func (s *UseCase) Clean(ctx context.Context, date time.Time) error {
	return s.st.Clear(ctx, date)
}

func (s *UseCase) storeStats(ctx context.Context, data interface{}) error {
	switch v := data.(type) {
	case *storage.CPUTimeStat:
		return s.st.StoreCPUTimeStat(ctx, *v)
	case *storage.LoadStat:
		return s.st.StoreLoadStat(ctx, *v)
	case *storage.NetworkStatesStat:
		return s.st.StoreNetworkStatesStat(ctx, *v)
	case []storage.NetworkStatsItem:
		return s.st.StoreNetworkStats(ctx, v)
	case []storage.UsageStatItem:
		return s.st.StoreUsageStats(ctx, v)
	case []storage.DiskIoStatItem:
		return s.st.StorDiskIoStats(ctx, v)
	case storage.ProtocolTalkerItem:
		return s.st.StoreProtocolTalkersStat(ctx, v)
	case storage.BpsItem:
		return s.st.StoreBpsTalkersStat(ctx, v)
	case nil:
		return nil
	default:
		return ErrCollector
	}
}

func (s *UseCase) collectPeriodic(ctx context.Context, duration time.Duration) {
	collectTicker := time.NewTicker(duration)
	for {
		select {
		case <-collectTicker.C:
			s.logger.Info("start collect periodic data")
			wg := sync.WaitGroup{}
			wg.Add(len(s.collectors))
			for n, c := range s.collectors {
				name := n
				go func(collector monitor.Collector) {
					defer wg.Done()
					item, err := collector.Grab(ctx)
					if err != nil {
						s.logger.Error(name+": failed to grab info", zap.Error(err))
					}
					err = s.storeStats(ctx, item)
					if err != nil {
						s.logger.Error(name+": failed to store info", zap.Error(err))
					}
				}(c)
			}
			wg.Wait()
			s.logger.Info("successful collect periodic data")
		case <-ctx.Done():
			s.logger.Info("data collection interrupted")
			return
		}
	}
}

func (s *UseCase) collectConstant(ctx context.Context) {
	s.logger.Info("start collect constant data")
	wg := sync.WaitGroup{}
	wg.Add(len(s.constantCollectors))
	for n, c := range s.constantCollectors {
		name := n
		go func(collector monitor.ConstantCollector) {
			defer wg.Done()
			stats, errC := collector.GrabStream(ctx)
			for {
				select {
				case stat, ok := <-stats:
					if !ok {
						return
					}
					err := s.storeStats(ctx, stat)
					if err != nil {
						s.logger.Error(name+": failed to store info", zap.Error(err))
					}
				case err, ok := <-errC:
					if !ok {
						continue
					}
					s.logger.Error(fmt.Sprintf("%s: error get content: %v", name, err), zap.Error(err))
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}
	wg.Wait()
}

func (s *UseCase) Collect(ctx context.Context, duration time.Duration) error {
	go func() {
		s.collectPeriodic(ctx, duration)
	}()
	go func() {
		s.collectConstant(ctx)
	}()
	return nil
}

func (s *UseCase) GetStats(ctx context.Context, duration int64) (Stats, error) {
	stats := Stats{}

	if s.isCPUEnable {
		lastCPUTimes, err := s.st.GetCPUTimeStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.CPUInfo = getAvgCPU(lastCPUTimes)
	}

	if s.isLoadEnable {
		lastLoadInfo, err := s.st.GetLoadStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.LoadInfo = getAvgLoad(lastLoadInfo)
	}

	if s.isDiskEnable {
		lastDiskIo, err := s.st.GetDiskIoStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.DiskIoInfo = getAvgDiskIo(lastDiskIo)

		lastDu, err := s.st.GetUsageStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.DiskUsageInfo = getUniqueDu(lastDu)
	}

	if s.isNetworkEnable {
		nsStats, err := s.st.GetNetworkStatesStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.NetworkStateInfo = getAvgNetworkStates(nsStats)

		nstStats, err := s.st.GetNetworkStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.NetworkStatisticsInfo = getUniqueNetworkStatistics(nstStats)
	}

	if s.isNetworkTalkersEnable {
		bpsStats, err := s.st.GetBpsTalkersStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.BpsTalkersInfo = getUniqueBps(bpsStats, duration)

		protocolStats, err := s.st.GetProtocolTalkersStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		stats.ProtocolTalkersInfo = getUniqueProtocolTalkers(protocolStats)
	}
	return stats, nil
}

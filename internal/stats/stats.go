package stats

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/cpu"
	diskusage "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/disk/usage"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/load"
	networkstates "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/network/states"
	networkstatistics "github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/network/statistics"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"go.uber.org/zap"
)

var ErrCollector = errors.New("unsupported collector type")

type Config struct {
	ScanDuration    time.Duration `mapstructure:"scan_duration"`
	CleanDuration   time.Duration `mapstructure:"clean_duration"`
	IsCPUEnable     bool          `mapstructure:"cpu_enable"`
	IsLoadEnable    bool          `mapstructure:"load_enable"`
	IsNetworkEnable bool          `mapstructure:"network_enable"`
	IsDiskEnable    bool          `mapstructure:"disk_enable"`
}

type Stats struct {
	CPUInfo               cpu.CPUTimeStat                 `json:"cpu_info"`                //nolint:tagliatelle
	LoadInfo              load.LoadStat                   `json:"load_info"`               //nolint:tagliatelle
	NetworkStateInfo      networkstates.NetworkStatesStat `json:"network_state_info"`      //nolint:tagliatelle
	NetworkStatisticsInfo networkstatistics.NetworkStats  `json:"network_statistics_info"` //nolint:tagliatelle
	DiskUsageInfo         diskusage.UsageStats            `json:"disk_usage_info"`         //nolint:tagliatelle
}

type StatsUseCase struct {
	logger        logger.Logger
	st            storage.Storage
	cleanDuration time.Duration
	collectors    []monitor.Collector
}

func New(cfg *Config, st storage.Storage, logger logger.Logger) StatsUseCase {
	collectors := make([]monitor.Collector, 0, 1)
	if cfg.IsCPUEnable {
		collectors = append(collectors, cpu.New(st))
	}
	if cfg.IsLoadEnable {
		collectors = append(collectors, load.New(st))
	}
	if cfg.IsNetworkEnable {
		collectors = append(collectors, networkstates.New(st))
		collectors = append(collectors, networkstatistics.New(st))
	}
	if cfg.IsDiskEnable {
		collectors = append(collectors, diskusage.New(st))
	}
	return StatsUseCase{
		collectors:    collectors,
		st:            st,
		cleanDuration: cfg.CleanDuration,
		logger:        logger,
	}
}

func (s *StatsUseCase) Clean(ctx context.Context) error {
	return s.st.Clear(ctx, time.Now().Add(-s.cleanDuration))
}

func (s *StatsUseCase) Collect(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(len(s.collectors))
	for _, c := range s.collectors {
		go func(collector monitor.Collector) {
			defer wg.Done()
			err := collector.Grab(ctx)
			if err != nil {
				s.logger.Error("failed to clear storage", zap.Error(err))
			}
		}(c)
	}
	wg.Wait()
	return nil
}

func (s *StatsUseCase) GetStats(ctx context.Context, duration int64) (Stats, error) {
	stats := Stats{}
	for _, collector := range s.collectors {
		statsItem, err := collector.GetStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		switch v := statsItem.(type) {
		case cpu.CPUTimeStat:
			stats.CPUInfo = v
		case load.LoadStat:
			stats.LoadInfo = v
		case networkstates.NetworkStatesStat:
			stats.NetworkStateInfo = v
		case networkstatistics.NetworkStats:
			stats.NetworkStatisticsInfo = v
		case diskusage.UsageStats:
			stats.DiskUsageInfo = v
		default:
			return stats, ErrCollector
		}
	}
	return stats, nil
}

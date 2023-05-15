package stats

import (
	"context"
	"errors"
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
	CPUInfo               storage.CPUTimeStat       `json:"cpu_info"`                //nolint:tagliatelle
	LoadInfo              storage.LoadStat          `json:"load_info"`               //nolint:tagliatelle
	NetworkStateInfo      storage.NetworkStatesStat `json:"network_state_info"`      //nolint:tagliatelle
	NetworkStatisticsInfo storage.NetworkStats      `json:"network_statistics_info"` //nolint:tagliatelle
	DiskUsageInfo         storage.UsageStats        `json:"disk_usage_info"`         //nolint:tagliatelle
	DiskIoInfo            storage.DiskIoStat        `json:"disk_io_info"`            //nolint:tagliatelle
}

type UseCase struct {
	logger             logger.Logger
	st                 storage.Storage
	cleanDuration      time.Duration
	collectors         []monitor.Collector
	constantCollectors []monitor.ConstantCollector
	storages           map[string]storage.Storage
}

func New(cfg *Config, st storage.Storage, logger logger.Logger) UseCase {
	collectors := make([]monitor.Collector, 0, 1)
	if cfg.IsCPUEnable {
		collectors = append(collectors, cpu.New(st, logger))
	}
	if cfg.IsLoadEnable {
		collectors = append(collectors, load.New(st, logger))
	}
	if cfg.IsNetworkEnable {
		collectors = append(collectors, networkstates.New(st, logger))
		collectors = append(collectors, networkstatistics.New(st, logger))
	}
	if cfg.IsDiskEnable {
		collectors = append(collectors, diskusage.New(st, logger))
		collectors = append(collectors, diskio.New(st, logger))
	}
	return UseCase{
		collectors:    collectors,
		st:            st,
		cleanDuration: cfg.CleanDuration,
		logger:        logger,
	}
}

func (s *UseCase) Clean(ctx context.Context) error {
	return s.st.Clear(ctx, time.Now().Add(-s.cleanDuration))
}

func (s *UseCase) collectPeriodic(ctx context.Context, duration time.Duration) {
	collectTicker := time.NewTicker(duration)
	for {
		select {
		case <-collectTicker.C:
			s.logger.Info("start collect periodic data")
			wg := sync.WaitGroup{}
			wg.Add(len(s.collectors))
			for _, c := range s.collectors {
				go func(collector monitor.Collector) {
					defer wg.Done()
					err := collector.Grab(ctx)
					if err != nil {
						s.logger.Error("failed to grab info", zap.Error(err))
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
	for _, c := range s.constantCollectors {
		go func(collector monitor.ConstantCollector) {
			defer wg.Done()
			err := collector.GrabSub(ctx)
			if err != nil {
				s.logger.Error("failed to grab info", zap.Error(err))
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
	for _, collector := range s.collectors {
		statsItem, err := collector.GetStats(ctx, duration)
		if err != nil {
			return stats, err
		}
		switch v := statsItem.(type) {
		case storage.CPUTimeStat:
			stats.CPUInfo = v
		case storage.LoadStat:
			stats.LoadInfo = v
		case storage.NetworkStatesStat:
			stats.NetworkStateInfo = v
		case storage.NetworkStats:
			stats.NetworkStatisticsInfo = v
		case storage.UsageStats:
			stats.DiskUsageInfo = v
		case storage.DiskIoStat:
			stats.DiskIoInfo = v
		default:
			return stats, ErrCollector
		}
	}
	return stats, nil
}

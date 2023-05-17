package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	memorystorage "github.com/BBeRsErKeRR/system-stats-monitor/internal/storage/memory"
	"go.uber.org/zap"
)

var ErrorScanPeriod = errors.New("error wait duration")

type Config struct {
	ScanDuration time.Duration `mapstructure:"scan_duration"`
	// CleanDuration time.Duration `mapstructure:"clean_duration"`
}

type App struct {
	logger       logger.Logger
	scanDuration time.Duration
	statsConfig  *stats.Config
}

func New(logger logger.Logger, config *Config, configS *stats.Config) *App {
	return &App{
		logger:       logger,
		scanDuration: config.ScanDuration,
		statsConfig:  configS,
	}
}

func (a *App) CreateStorage() storage.Storage {
	return memorystorage.New()
}

func (a *App) CreateUseCase(st storage.Storage) stats.UseCase {
	return stats.New(a.statsConfig, st, a.logger)
}

func (a *App) StartMonitoring(ctx context.Context, rd, wp int64, u stats.UseCase) (<-chan stats.Stats, error) {
	res := make(chan stats.Stats)
	responseTicker := time.NewTicker(time.Duration(rd) * time.Second)
	waitDuration := time.Duration(wp) * time.Second
	cleanTicker := time.NewTicker(waitDuration)

	if waitDuration <= a.scanDuration {
		return res, ErrorScanPeriod
	}

	go func() {
		a.logger.Info(fmt.Sprintf("collect data from period: %s", a.scanDuration))
		if err := u.Collect(ctx, a.scanDuration); err != nil {
			a.logger.Error("failed to collect metrics: " + err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-cleanTicker.C:
				err := u.Clean(ctx, time.Now().Add(-waitDuration))
				a.logger.Info("clean done")
				if err != nil {
					a.logger.Error("failed to clear storage: ", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(res)
		startTime := time.Now().Add(waitDuration)
		for {
			select {
			case <-responseTicker.C:
				if rd < wp && time.Now().Before(startTime) {
					continue
				}
				stats, err := u.GetStats(ctx, wp)
				if err != nil {
					a.logger.Error("fail get stats", zap.Error(err))
				}
				res <- stats
			case <-ctx.Done():
				a.logger.Info("sending data interrupted")
				return
			}
		}
	}()

	return res, nil
}

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	memorystorage "github.com/BBeRsErKeRR/system-stats-monitor/internal/storage/memory"
	"go.uber.org/zap"
)

type App struct {
	logger logger.Logger
	u      stats.UseCase
}

func New(logger logger.Logger, u stats.UseCase) *App {
	return &App{
		logger: logger,
		u:      u,
	}
}

func (a *App) CollectData(ctx context.Context, duration time.Duration) error {
	a.logger.Info(fmt.Sprintf("collect data from period: %s", duration))
	err := a.u.Collect(ctx, duration)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) StartMonitoring(ctx context.Context, respDuration, statsDuration int64) (<-chan stats.Stats, error) {
	res := make(chan stats.Stats)
	responseTicker := time.NewTicker(time.Duration(respDuration) * time.Second)

	go func() {
		defer close(res)
		for {
			select {
			case <-responseTicker.C:
				stats, err := a.u.GetStats(ctx, statsDuration)
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

func GetStorage() storage.Storage {
	return memorystorage.New()
}

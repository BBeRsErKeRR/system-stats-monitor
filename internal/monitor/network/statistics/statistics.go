package networkstatistics

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
)

type StatCollector struct {
	logger logger.Logger
}

func New(logger logger.Logger) *StatCollector {
	return &StatCollector{
		logger: logger,
	}
}

func (c *StatCollector) Grab(ctx context.Context) (interface{}, error) {
	c.logger.Info("start collect data")
	stats, err := getNS(ctx)
	if err != nil {
		return nil, err
	}
	c.logger.Info("successful collect data")
	return stats, nil
}

func (c *StatCollector) CheckExecution(ctx context.Context) error {
	return checkCall(ctx)
}

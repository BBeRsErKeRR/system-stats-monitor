package cpu

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

var commandRunner = getCPUTimes

func (c *StatCollector) Grab(ctx context.Context) (interface{}, error) {
	c.logger.Info("start collect data")
	times, err := commandRunner(ctx)
	if err != nil {
		return nil, err
	}
	c.logger.Info("successful collect data")
	return times, nil
}

func (s *StatCollector) CheckCall(ctx context.Context) error {
	return checkCall(ctx)
}

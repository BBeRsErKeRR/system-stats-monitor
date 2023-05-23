package bpstalkers

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

func (c *StatCollector) GrabSub(ctx context.Context) (<-chan interface{}, <-chan error) {
	return getBps(ctx)
}

func (s *StatCollector) CheckCall(ctx context.Context) error {
	return checkCall(ctx)
}

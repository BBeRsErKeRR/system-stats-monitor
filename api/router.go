package router

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
)

type Application interface {
	StartMonitoring(context.Context, int64, int64) (<-chan stats.Stats, error)
}

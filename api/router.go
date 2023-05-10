package router

import (
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
)

type Application interface {
	StartMonitoring(int64, int64) (<-chan monitor.Stats, error)
}

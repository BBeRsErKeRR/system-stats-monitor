//go:build darwin
// +build darwin

package cpu

import (
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getCPUTimes() (*storage.CPUTimeStat, error) {
	return nil, monitor.ErrNotImplemented
}

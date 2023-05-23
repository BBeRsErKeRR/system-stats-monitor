//go:build darwin
// +build darwin

package cpu

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getCPUTimes(ctx context.Context) (*storage.CPUTimeStat, error) {
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}

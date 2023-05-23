//go:build windows
// +build windows

package diskusage

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getDU(ctx context.Context) ([]storage.UsageStatItem, error) {
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}

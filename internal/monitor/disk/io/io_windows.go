//go:build windows
// +build windows

package diskio

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func collectDiskIo(ctx context.Context) ([]storage.DiskIoStatItem, error) { //nolint:revive
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error { //nolint:revive
	return monitor.ErrNotImplemented
}

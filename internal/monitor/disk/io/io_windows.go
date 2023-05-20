//go:build windows
// +build windows

package diskio

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func collectDiskIo(ctx context.Context) ([]interface{}, error) {
	return nil, monitor.ErrNotImplemented
}

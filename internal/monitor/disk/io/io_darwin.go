//go:build darwin
// +build darwin

package diskio

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
)

func collectDiskIo(ctx context.Context) ([]interface{}, error) {
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}

//go:build windows
// +build windows

package networkstatistics

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
)

func getNS(ctx context.Context) ([]interface{}, error) {
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}

//go:build windows
// +build windows

package networkstatistics

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
)

func getNS(ctx context.Context) ([]interface{}, error) { //nolint:revive
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error { //nolint:revive
	return monitor.ErrNotImplemented
}

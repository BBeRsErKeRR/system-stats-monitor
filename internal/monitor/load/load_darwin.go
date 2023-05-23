//go:build darwin
// +build darwin

package load

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getLoad() (*storage.LoadStat, error) {
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error { //nolint:revive
	return monitor.ErrNotImplemented
}

//go:build darwin
// +build darwin

package networkstates

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getNS(ctx context.Context) (*storage.NetworkStatesStat, error) {
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}

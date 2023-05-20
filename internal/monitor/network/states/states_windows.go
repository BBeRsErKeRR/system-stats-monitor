//go:build windows
// +build windows

package networkstates

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getNS(ctx context.Context) (*storage.NetworkStatesStat, error) {
	return nil, monitor.ErrNotImplemented
}

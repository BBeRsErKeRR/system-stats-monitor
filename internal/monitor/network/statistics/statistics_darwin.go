//go:build darwin
// +build darwin

package networkstatistics

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
)

func getNS(ctx context.Context) ([]interface{}, error) {
	return nil, monitor.ErrNotImplemented
}

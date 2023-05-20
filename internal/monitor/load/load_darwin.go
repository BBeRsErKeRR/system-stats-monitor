//go:build darwin
// +build darwin

package load

import (
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getLoad() (*storage.LoadStat, error) {
	return nil, monitor.ErrNotImplemented
}

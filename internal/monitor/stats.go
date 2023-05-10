package monitor

import (
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor/cpu"
)

type Stats struct {
	CPUInfo cpu.CPUTimeStat `json:"cpu_info"` //nolint:tagliatelle
}

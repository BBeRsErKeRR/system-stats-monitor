//go:build linux
// +build linux

package load

import (
	"context"
	"strconv"
	"strings"

	files "github.com/BBeRsErKeRR/system-stats-monitor/pkg/files"
)

func parseStatLine(line string) (*LoadStat, error) {
	values := strings.Fields(line)

	load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}
	load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}
	load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	ret := &LoadStat{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}
	return ret, nil
}

func getLoad(ctx context.Context) (*LoadStat, error) {
	lines, err := files.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	return parseStatLine(lines)
}

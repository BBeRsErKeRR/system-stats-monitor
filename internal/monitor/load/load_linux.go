//go:build linux
// +build linux

package load

import (
	"context"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	files "github.com/BBeRsErKeRR/system-stats-monitor/pkg/files"
)

func parseStatLine(line string) (*storage.LoadStat, error) {
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

	ret := &storage.LoadStat{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}
	return ret, nil
}

func getLoad() (*storage.LoadStat, error) {
	lines, err := files.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	return parseStatLine(lines)
}

func checkCall(ctx context.Context) error {
	_, err := files.ReadLinesOffsetN("/proc/loadavg", 0, 1)
	return err
}

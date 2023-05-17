//go:build linux
// +build linux

package diskusage

import (
	"context"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseDfOut(outK, outI string) ([]storage.UsageStatItem, error) {
	lines := strings.Split(outK, "\n")
	linesI := strings.Split(outI, "\n")
	result := make([]storage.UsageStatItem, 0, len(outK)-1)
	buff := make(map[string]storage.UsageStatItem)

	for _, line := range lines[1:] {
		values := strings.Fields(line)
		if len(values) < 1 {
			continue
		}
		used, err := strconv.ParseInt(values[2], 10, 32)
		if err != nil {
			return nil, err
		}

		available, err := strconv.ParseInt(values[3], 10, 32)
		if err != nil {
			return nil, err
		}

		item := storage.UsageStatItem{
			Path:             values[5],
			Fstype:           values[0],
			Used:             used / 1024,
			AvailablePercent: (float64(available) / float64(available+used)) * 100.0,
		}
		buff[values[5]] = item
	}

	for _, line := range linesI[1:] {
		values := strings.Fields(line)
		if len(values) < 1 {
			continue
		}
		used, err := strconv.ParseInt(values[2], 10, 32)
		if err != nil {
			return nil, err
		}

		available, err := strconv.ParseInt(values[3], 10, 32)
		if err != nil {
			return nil, err
		}

		item := buff[values[5]]
		item.InodesUsed = used / 1024
		item.InodesAvailablePercent = (float64(available) / float64(available+used)) * 100.0
		result = append(result, item)
	}
	return result, nil
}

func getDU(ctx context.Context) ([]storage.UsageStatItem, error) {
	outK, err := command.WithContext(ctx, "df", "-k")
	if err != nil {
		return nil, err
	}
	outI, err := command.WithContext(ctx, "df", "-i")
	if err != nil {
		return nil, err
	}
	return parseDfOut(string(outK), string(outI))
}

//go:build linux
// +build linux

package diskusage

import (
	"context"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseDfOut(outK, outI string) ([]interface{}, error) {
	lines := strings.Split(outK, "\n")
	linesI := strings.Split(outI, "\n")
	result := make([]interface{}, 0, len(outK)-1)
	buff := make(map[string]UsageStatItem)

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

		item := UsageStatItem{
			Path:             values[5],
			Fstype:           values[0],
			Used:             int64(used) / 1024,
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
		item.InodesUsed = int64(used) / 1024
		item.InodesAvailablePercent = (float64(available) / float64(available+used)) * 100.0
		result = append(result, item)
	}
	return result, nil
}

func GetDU(ctx context.Context) ([]interface{}, error) {
	outK, err := command.CommandWithContext(ctx, "df", "-k")
	if err != nil {
		return nil, err
	}
	outI, err := command.CommandWithContext(ctx, "df", "-i")
	if err != nil {
		return nil, err
	}
	return parseDfOut(string(outK), string(outI))
}

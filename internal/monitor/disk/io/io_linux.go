//go:build linux
// +build linux

package diskio

import (
	"context"
	"strconv"
	"strings"

	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseSSOut(output string) ([]interface{}, error) {
	lines := strings.Split(output, "\n")
	result := make([]interface{}, 0, len(output)-3)
	for _, line := range lines[3:] {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		tps, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}

		kb_read_s, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}

		kb_write_s, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return nil, err
		}

		result = append(result, NewDiskIoStatItem(fields[0], tps, kb_read_s, kb_write_s))
	}

	return result, nil
}

func collectDiskIo(ctx context.Context) ([]interface{}, error) {
	out, err := command.CommandWithContext(ctx, "iostat", "-d", "-k")
	if err != nil {
		return nil, err
	}
	return parseSSOut(string(out))
}

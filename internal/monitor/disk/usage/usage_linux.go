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

type DuDTO struct {
	Path      string
	Fstype    string
	Used      int64
	Available int64
}

func parseFields(line string) (DuDTO, error) {
	var res DuDTO
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return res, nil
	}
	filesystemMaxLength := len(fields) - 6
	filesystem := strings.Join(fields[:filesystemMaxLength+1], " ")

	used, err := strconv.ParseInt(fields[filesystemMaxLength+2], 10, 32)
	if err != nil {
		return res, err
	}

	if used < 0 {
		used = 0
	}

	available, err := strconv.ParseInt(fields[filesystemMaxLength+3], 10, 32)
	if err != nil {
		return res, err
	}

	if available < 0 {
		available = 0
	}

	mount := fields[filesystemMaxLength+5]

	return DuDTO{
		Path:      mount,
		Fstype:    filesystem,
		Used:      used,
		Available: available,
	}, nil
}

func parseDfOut(outK, outI string) ([]storage.UsageStatItem, error) {
	lines := strings.Split(outK, "\n")
	linesI := strings.Split(outI, "\n")
	result := make([]storage.UsageStatItem, 0, len(lines)-1)
	buff := map[string]storage.UsageStatItem{}
	for _, line := range lines[1:] {
		dto, err := parseFields(line)
		if err != nil {
			return nil, err
		}
		if dto == (DuDTO{}) {
			continue
		}

		item := storage.UsageStatItem{
			Path:   dto.Path,
			Fstype: dto.Fstype,
			Used:   dto.Used / 1024,
		}
		if dto.Available == 0 && dto.Used == 0 {
			item.AvailablePercent = 0.00
		} else {
			item.AvailablePercent = (float64(dto.Available) / float64(dto.Available+dto.Used)) * 100.0
		}
		buff[dto.Path] = item
	}

	for _, line := range linesI[1:] {
		dto, err := parseFields(line)
		if err != nil {
			return nil, err
		}
		if dto == (DuDTO{}) {
			continue
		}
		item := buff[dto.Path]
		item.InodesUsed = dto.Used / 1024
		if dto.Available == 0 && dto.Used == 0 {
			item.InodesAvailablePercent = 0.00
		} else {
			item.InodesAvailablePercent = (float64(dto.Available) / float64(dto.Available+dto.Used)) * 100.0
		}
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

func checkCall(ctx context.Context) error {
	_, err := command.WithContext(ctx, "df")
	return err
}

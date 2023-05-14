package diskusage

import (
	"context"
	"fmt"
	"sort"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type UsageStatItem struct {
	Path                   string  `json:"path"`
	Fstype                 string  `json:"fstype"`
	Used                   int64   `json:"used"`
	AvailablePercent       float64 `json:"available_percent"`
	InodesUsed             int64   `json:"inodes_used"`
	InodesAvailablePercent float64 `json:"inodes_available_percent"`
}

type UsageStats struct {
	Items []UsageStatItem
}

type UsageCollector struct {
	name string
	st   storage.Storage
}

func New(st storage.Storage) *UsageCollector {
	return &UsageCollector{
		name: "du",
		st:   st,
	}
}

func (c *UsageCollector) Grab(ctx context.Context) error {
	stats, err := getDU(ctx)
	if err != nil {
		return err
	}
	for _, stat := range stats {
		err = c.st.StoreStats(ctx, c.name, stat)
		if err != nil {
			return err
		}
	}
	return nil
}

func (as *UsageCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	stats, err := as.st.GetStats(ctx, as.name, period)
	if err != nil {
		return nil, err
	}
	return UsageStats{
		Items: unique(stats),
	}, nil
}

func unique(intSlice []storage.Metric) []UsageStatItem {
	keys := make(map[string]bool)
	list := make([]UsageStatItem, 0, len(intSlice))
	sort.Slice(intSlice, func(i, j int) bool {
		return intSlice[i].Date.Before(intSlice[j].Date)
	})

	for _, fact := range intSlice {
		stat := fact.StatInfo.(UsageStatItem)
		entry := fmt.Sprintf("%v/%v", stat.Path, stat.Fstype)
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, stat)
		}
	}
	return list
}

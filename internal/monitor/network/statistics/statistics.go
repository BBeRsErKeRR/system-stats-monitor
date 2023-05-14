package networkstatistics

import (
	"context"
	"fmt"
	"sort"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type NetworkStatsItem struct {
	Command  string
	PID      int32
	User     int32
	Protocol string
	Port     int32
}

type NetworkStats struct {
	Items []NetworkStatsItem
}

type NSCollector struct {
	name string
	st   storage.Storage
}

func New(st storage.Storage) *NSCollector {
	return &NSCollector{
		name: "network_statistics",
		st:   st,
	}
}

func (c *NSCollector) Grab(ctx context.Context) error {
	stats, err := getNS(ctx)
	if err != nil {
		return err
	}
	return c.st.BulkStoreStats(ctx, c.name, stats)
}

func (as *NSCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	nsStats, err := as.st.GetStats(ctx, as.name, period)
	if err != nil {
		return nil, err
	}
	return NetworkStats{
		Items: unique(nsStats),
	}, nil
}

func unique(intSlice []storage.Metric) []NetworkStatsItem {
	keys := make(map[string]bool)
	list := make([]NetworkStatsItem, 0, len(intSlice))

	sort.Slice(intSlice, func(i, j int) bool {
		return intSlice[i].Date.Before(intSlice[j].Date)
	})

	for _, fact := range intSlice {
		stat := fact.StatInfo.(NetworkStatsItem)
		entry := fmt.Sprintf("%v/%v/%v/%v", stat.Command, stat.Protocol, stat.PID, stat.Port)
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, stat)
		}
	}
	return list
}

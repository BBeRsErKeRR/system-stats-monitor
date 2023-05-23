package networkstatistics

import (
	"context"
	"fmt"
	"sort"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type StatCollector struct {
	st     storage.Storage
	logger logger.Logger
}

func New(st storage.Storage, logger logger.Logger) *StatCollector {
	return &StatCollector{
		st:     st,
		logger: logger,
	}
}

func (c *StatCollector) Grab(ctx context.Context) error {
	c.logger.Info("start collect data")
	stats, err := getNS(ctx)
	if err != nil {
		return err
	}
	err = c.st.StoreNetworkStats(ctx, stats)
	if err != nil {
		return err
	}
	c.logger.Info("successful collect data")
	return nil
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	nsStats, err := c.st.GetNetworkStats(ctx, period)
	if err != nil {
		return nil, err
	}
	return storage.NetworkStats{
		Items: unique(nsStats),
	}, nil
}

func unique(intSlice []storage.Metric) []storage.NetworkStatsItem {
	keys := make(map[string]bool)
	list := make([]storage.NetworkStatsItem, 0, len(intSlice))

	sort.Slice(intSlice, func(i, j int) bool {
		return intSlice[i].Date.Before(intSlice[j].Date)
	})

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.NetworkStatsItem)
		entry := fmt.Sprintf("%v/%v/%v/%v", stat.Command, stat.Protocol, stat.PID, stat.Port)
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, stat)
		}
	}
	return list
}

func (s *StatCollector) CheckCall(ctx context.Context) error {
	return checkCall(ctx)
}

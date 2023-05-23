package diskusage

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
	stats, err := getDU(ctx)
	if err != nil {
		return err
	}
	for _, stat := range stats {
		err = c.st.StoreUsageStat(ctx, stat)
		if err != nil {
			return err
		}
	}
	c.logger.Info("successful collect data")
	return nil
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	stats, err := c.st.GetUsageStats(ctx, period)
	if err != nil {
		return nil, err
	}
	return storage.UsageStats{
		Items: unique(stats),
	}, nil
}

func unique(intSlice []storage.Metric) []storage.UsageStatItem {
	keys := make(map[string]bool)
	list := make([]storage.UsageStatItem, 0, len(intSlice))
	sort.Slice(intSlice, func(i, j int) bool {
		return intSlice[i].Date.Before(intSlice[j].Date)
	})

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.UsageStatItem)
		entry := fmt.Sprintf("%v/%v", stat.Path, stat.Fstype)
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

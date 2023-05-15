package diskusage

import (
	"context"
	"fmt"
	"sort"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type UsageCollector struct {
	name   string
	st     storage.Storage
	logger logger.Logger
}

func New(st storage.Storage, logger logger.Logger) *UsageCollector {
	return &UsageCollector{
		name:   "du",
		st:     st,
		logger: logger,
	}
}

func (c *UsageCollector) Grab(ctx context.Context) error {
	c.logger.Info("start collect data")
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
	c.logger.Info("successful collect data")
	return nil
}

func (c *UsageCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	stats, err := c.st.GetStats(ctx, c.name, period)
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

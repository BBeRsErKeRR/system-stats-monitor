package networkstates

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type StatCollector struct {
	name   string
	st     storage.Storage
	logger logger.Logger
}

func New(st storage.Storage, logger logger.Logger) *StatCollector {
	return &StatCollector{
		name:   "network_sates",
		st:     st,
		logger: logger,
	}
}

func (c *StatCollector) Grab(ctx context.Context) error {
	c.logger.Info("start collect data")
	stat, err := getNS(ctx)
	if err != nil {
		return err
	}
	err = c.st.StoreStats(ctx, c.name, *stat)
	if err != nil {
		return err
	}
	c.logger.Info("successful collect data")
	return nil
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	avgNs := make(map[string]int32)
	nsStats, err := c.st.GetStats(ctx, c.name, period)
	if err != nil {
		return nil, err
	}
	for _, fact := range nsStats {
		stat := fact.StatInfo.(storage.NetworkStatesStat)
		for name, counter := range stat.Counters {
			_, ok := avgNs[name]
			if ok {
				avgNs[name] += counter
			} else {
				avgNs[name] = counter
			}
		}
	}
	lengthStat := int32(len(nsStats))
	for name, counter := range avgNs {
		avgNs[name] = counter / lengthStat
	}
	return storage.NetworkStatesStat{
		Counters: avgNs,
	}, nil
}

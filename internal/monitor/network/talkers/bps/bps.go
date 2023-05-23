package bpstalkers

import (
	"context"
	"fmt"
	"sort"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"go.uber.org/zap"
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

func (c *StatCollector) GrabSub(ctx context.Context) error {
	c.logger.Info("start collect data")
	stats, errC := getBps(ctx)

	for {
		select {
		case stat, ok := <-stats:
			if !ok {
				return nil
			}
			err := c.st.StoreBpsTalkersStat(ctx, stat)
			if err != nil {
				return err
			}
		case err, ok := <-errC:
			if !ok {
				continue
			}
			c.logger.Error(fmt.Sprintf("error get content: %v", err), zap.Error(err))
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	nsStats, err := c.st.GetBpsTalkersStats(ctx, period)
	if err != nil {
		return nil, err
	}
	return storage.BpsTalkersStats{
		Items: c.collectUnique(nsStats, period),
	}, nil
}

func (c *StatCollector) collectUnique(intSlice []storage.Metric, period int64) []storage.BpsItem {
	c.logger.Info("get protocols")
	unique := make(map[string]storage.BpsItem)
	list := make([]storage.BpsItem, 0, len(intSlice))

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.BpsItem)
		entry := fmt.Sprintf("%s->%s", stat.Source, stat.Destination)
		item, value := unique[entry]
		if !value {
			unique[entry] = stat
		} else {
			item.Numbers += stat.Numbers
			unique[entry] = item
		}
	}
	seconds := float64(period)
	for _, elem := range unique {
		elem.Bps = elem.Numbers / seconds
		list = append(list, elem)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Bps > list[j].Bps
	})

	return list
}

func (s *StatCollector) CheckCall(ctx context.Context) error {
	return checkCall(ctx)
}

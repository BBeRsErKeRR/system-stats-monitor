package protocoltalkers

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
	stats, errC := getTalkers(ctx)

	for {
		select {
		case stat, ok := <-stats:
			if !ok {
				return nil
			}
			err := c.st.StoreProtocolTalkersStat(ctx, stat)
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
	nsStats, err := c.st.GetProtocolTalkersStats(ctx, period)
	if err != nil {
		return nil, err
	}
	return storage.ProtocolTalkersStats{
		Items: c.collectUnique(nsStats),
	}, nil
}

func (c *StatCollector) collectUnique(intSlice []storage.Metric) []storage.ProtocolTalkerItem {
	var sumBytes float64
	c.logger.Info("get protocols")
	unique := make(map[string]storage.ProtocolTalkerItem)
	list := make([]storage.ProtocolTalkerItem, 0, len(intSlice))

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.ProtocolTalkerItem)
		entry := stat.Protocol
		item, value := unique[entry]
		if !value {
			unique[entry] = stat
		} else {
			item.SendBytes += stat.SendBytes
			unique[entry] = item
		}
		sumBytes += stat.SendBytes
	}

	for _, elem := range unique {
		elem.BytesPercentage = (elem.SendBytes / sumBytes) * 100.0
		list = append(list, elem)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Protocol < list[j].Protocol
	})

	return list
}

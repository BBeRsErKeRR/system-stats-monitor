package networkstates

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type NetworkStatesStat struct {
	Counters map[string]int32
}

type NSCollector struct {
	name string
	st   storage.Storage
}

func New(st storage.Storage) *NSCollector {
	return &NSCollector{
		name: "network_sates",
		st:   st,
	}
}

func (c *NSCollector) Grab(ctx context.Context) error {
	stat, err := getNS(ctx)
	if err != nil {
		return err
	}
	return c.st.StoreStats(ctx, c.name, *stat)
}

func (as *NSCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	avgNs := make(map[string]int32)
	nsStats, err := as.st.GetStats(ctx, as.name, period)
	if err != nil {
		return nil, err
	}
	for _, fact := range nsStats {
		stat := fact.StatInfo.(NetworkStatesStat)
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
	return NetworkStatesStat{
		Counters: avgNs,
	}, nil
}

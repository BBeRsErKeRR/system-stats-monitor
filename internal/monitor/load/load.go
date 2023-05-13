package load

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type LoadStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type LoadStatCollector struct {
	name string
	st   storage.Storage
}

func New(st storage.Storage) *LoadStatCollector {
	return &LoadStatCollector{
		name: "load",
		st:   st,
	}
}

func (c *LoadStatCollector) Grab(ctx context.Context) error {
	la, err := GetLoad(ctx)
	if err != nil {
		return err
	}
	return c.st.StoreStats(ctx, c.name, *la)
}

func (as *LoadStatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	var sumLoad1, sumLoad5, sumLoad15 float64
	lastLoadStats, err := as.st.GetStats(ctx, as.name, period)
	if err != nil {
		return nil, err
	}
	for _, metric := range lastLoadStats {
		stat := metric.StatInfo.(LoadStat)
		sumLoad1 += stat.Load1
		sumLoad5 += stat.Load5
		sumLoad15 += stat.Load15
	}
	totalLen := len(lastLoadStats)
	return LoadStat{
		Load1:  sumLoad1 / float64(totalLen),
		Load5:  sumLoad5 / float64(totalLen),
		Load15: sumLoad15 / float64(totalLen),
	}, nil
}

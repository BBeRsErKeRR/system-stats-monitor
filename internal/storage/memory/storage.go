package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	metrics map[string]storage.Metric
	sync.RWMutex
}

func (st *Storage) StoreStats(ctx context.Context, name string, data interface{}) error {
	st.Lock()
	defer st.Unlock()
	id := uuid.New().String()
	st.metrics[id] = storage.CreateMetric(name, data)
	return nil
}

func (st *Storage) GetStats(ctx context.Context, name string, period int64) ([]storage.Metric, error) {
	end := time.Now()
	start := time.Now().Add(-time.Duration(period) * time.Second)
	st.RLock()
	defer st.RUnlock()
	res := make([]storage.Metric, 0, len(st.metrics))
	for _, e := range st.metrics {
		if e.Name == name && e.Date.After(start) && e.Date.Before(end) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (st *Storage) Clear(ctx context.Context, date time.Time) error {
	st.Lock()
	defer st.Unlock()
	for id, m := range st.metrics {
		if m.Date.After(date) {
			delete(st.metrics, id)
		}
	}
	return nil
}

func (st *Storage) Connect(ctx context.Context) error {
	return nil
}

func (st *Storage) Close(ctx context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{
		metrics: map[string]storage.Metric{},
	}
}

package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/google/uuid"
)

type MemStore struct {
	st map[string]storage.Metric
	sync.RWMutex
}

func (ms *MemStore) GetStats(ctx context.Context, period int64) ([]storage.Metric, error) { //nolint:revive
	end := time.Now()
	start := time.Now().Add(-time.Duration(period) * time.Second)
	ms.RLock()
	defer ms.RUnlock()
	res := make([]storage.Metric, 0, len(ms.st))
	for _, e := range ms.st {
		if e.Date.After(start) && e.Date.Before(end) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (ms *MemStore) StoreStats(ctx context.Context, data interface{}) error { //nolint:revive
	ms.Lock()
	defer ms.Unlock()
	id := uuid.New().String()
	ms.st[id] = storage.CreateMetric(data, time.Now())
	return nil
}

func bulkStore[T any](ms *MemStore, data []T) error {
	ms.Lock()
	defer ms.Unlock()
	date := time.Now()
	for _, item := range data {
		id := uuid.New().String()
		ms.st[id] = storage.CreateMetric(item, date)
	}
	return nil
}

func (ms *MemStore) Clear(ctx context.Context, date time.Time) error { //nolint:revive
	ms.Lock()
	defer ms.Unlock()
	for id, m := range ms.st {
		if m.Date.Before(date) {
			delete(ms.st, id)
		}
	}
	return nil
}

func NewMemStore() MemStore {
	return MemStore{
		st: map[string]storage.Metric{},
	}
}

type Storage struct {
	cpuSt       MemStore
	loadSt      MemStore
	nsSt        MemStore
	nStatsSt    MemStore
	duSt        MemStore
	dioSt       MemStore
	protoTalkSt MemStore
	bpsTalkSt   MemStore
}

func (st *Storage) Clear(ctx context.Context, date time.Time) error {
	var errC chan error
	wg := sync.WaitGroup{}
	wg.Add(8)
	go func() {
		defer wg.Done()
		err := st.cpuSt.Clear(ctx, date)
		if err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.loadSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.nsSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.nStatsSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.duSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.dioSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.protoTalkSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := st.bpsTalkSt.Clear(ctx, date); err != nil {
			errC <- err
		}
	}()
	wg.Wait()
	select {
	case err := <-errC:
		return err
	default:
		return nil
	}
}

func (st *Storage) Connect(ctx context.Context) error { //nolint:revive
	return nil
}

func (st *Storage) Close(ctx context.Context) error { //nolint:revive
	return nil
}

func (st *Storage) StoreCPUTimeStat(ctx context.Context, data storage.CPUTimeStat) error {
	return st.cpuSt.StoreStats(ctx, data)
}

func (st *Storage) GetCPUTimeStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.cpuSt.GetStats(ctx, period)
}

func (st *Storage) StoreLoadStat(ctx context.Context, data storage.LoadStat) error {
	return st.loadSt.StoreStats(ctx, data)
}

func (st *Storage) GetLoadStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.loadSt.GetStats(ctx, period)
}

func (st *Storage) StoreNetworkStatesStat(ctx context.Context, data storage.NetworkStatesStat) error {
	return st.nsSt.StoreStats(ctx, data)
}

func (st *Storage) GetNetworkStatesStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.nsSt.GetStats(ctx, period)
}

func (st *Storage) StoreNetworkStats(ctx context.Context, data []storage.NetworkStatsItem) error { //nolint:revive
	return bulkStore(&st.nStatsSt, data)
}

func (st *Storage) GetNetworkStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.nStatsSt.GetStats(ctx, period)
}

func (st *Storage) StoreUsageStats(ctx context.Context, data []storage.UsageStatItem) error { //nolint:revive
	return bulkStore(&st.duSt, data)
}

func (st *Storage) GetUsageStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.duSt.GetStats(ctx, period)
}

func (st *Storage) StorDiskIoStats(ctx context.Context, data []storage.DiskIoStatItem) error { //nolint:revive
	return bulkStore(&st.dioSt, data)
}

func (st *Storage) GetDiskIoStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.dioSt.GetStats(ctx, period)
}

func (st *Storage) StoreProtocolTalkersStat(ctx context.Context, data storage.ProtocolTalkerItem) error {
	return st.protoTalkSt.StoreStats(ctx, data)
}

func (st *Storage) GetProtocolTalkersStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.protoTalkSt.GetStats(ctx, period)
}

func (st *Storage) StoreBpsTalkersStat(ctx context.Context, data storage.BpsItem) error {
	return st.bpsTalkSt.StoreStats(ctx, data)
}

func (st *Storage) GetBpsTalkersStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return st.bpsTalkSt.GetStats(ctx, period)
}

func New() *Storage {
	return &Storage{
		cpuSt:       NewMemStore(),
		loadSt:      NewMemStore(),
		nsSt:        NewMemStore(),
		nStatsSt:    NewMemStore(),
		duSt:        NewMemStore(),
		dioSt:       NewMemStore(),
		protoTalkSt: NewMemStore(),
		bpsTalkSt:   NewMemStore(),
	}
}

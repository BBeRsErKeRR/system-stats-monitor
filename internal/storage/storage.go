package storage

import (
	"context"
	"time"
)

type Metric struct {
	Name     string
	Date     time.Time
	StatInfo interface{}
}

func CreateMetric(name string, data interface{}, date time.Time) Metric {
	return Metric{
		Name:     name,
		Date:     date,
		StatInfo: data,
	}
}

type Storage interface {
	StoreStats(context.Context, string, interface{}) error
	BulkStoreStats(context.Context, string, []interface{}) error
	GetStats(context.Context, string, int64) ([]Metric, error)
	Clear(context.Context, time.Time) error
	Connect(context.Context) error
	Close(context.Context) error
}

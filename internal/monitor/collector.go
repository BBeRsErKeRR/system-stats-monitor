package monitor

import "context"

type Collector interface {
	Grab(context.Context) error
	GetStats(context.Context, int64) (interface{}, error)
}

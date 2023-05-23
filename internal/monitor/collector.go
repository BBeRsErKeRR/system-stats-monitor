package monitor

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("collector not implemented")

type Collector interface {
	Grab(context.Context) error
	GetStats(context.Context, int64) (interface{}, error)
	CheckCall(context.Context) error
}

type ConstantCollector interface {
	GrabSub(context.Context) error
	GetStats(context.Context, int64) (interface{}, error)
	CheckCall(context.Context) error
}

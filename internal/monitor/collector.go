package monitor

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("collector not implemented")

type Collector interface {
	Grab(context.Context) (interface{}, error)
	CheckCall(context.Context) error
}

type ConstantCollector interface {
	GrabSub(context.Context) (<-chan interface{}, <-chan error)
	CheckCall(context.Context) error
}

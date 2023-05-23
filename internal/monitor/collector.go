package monitor

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("collector not implemented")

type Collector interface {
	Grab(context.Context) (interface{}, error)
	CheckExecution(context.Context) error
}

type StreamCollector interface {
	GrabSub(context.Context) (<-chan interface{}, <-chan error)
	CheckExecution(context.Context) error
}

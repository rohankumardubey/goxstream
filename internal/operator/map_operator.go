package operator

import "goxstream/internal/model"

type MapFunc func(event model.Event) model.Event

type MapOperator struct {
    name string
    fn   MapFunc
}

func NewMapOperator(name string, fn MapFunc) *MapOperator {
    return &MapOperator{name: name, fn: fn}
}

func (op *MapOperator) Name() string { return op.name }

func (op *MapOperator) Process(event model.Event) []model.Event {
    return []model.Event{op.fn(event)}
}

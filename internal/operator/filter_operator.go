package operator

import (
    "goxstream/internal/model"
)

// Define a filter function type
type FilterFunc func(event model.Event) bool

type FilterOperator struct {
    name string
    fn   FilterFunc
}

func NewFilterOperator(name string, fn FilterFunc) *FilterOperator {
    return &FilterOperator{name: name, fn: fn}
}

func (op *FilterOperator) Name() string { return op.name }

// Only emit events for which fn(event) == true
func (op *FilterOperator) Process(event model.Event) []model.Event {
    if op.fn(event) {
        return []model.Event{event}
    }
    return nil
}

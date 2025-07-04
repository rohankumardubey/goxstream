package operator

import (
    "goxstream/internal/model"
)

type TumblingWindowOperator struct {
    name  string
    size  int
	buffer []model.Event
    inner Operator // e.g., a reduce operator
}

func NewTumblingWindowOperator(name string, size int, inner Operator) *TumblingWindowOperator {
    return &TumblingWindowOperator{name: name, size: size, inner: inner}
}

func (op *TumblingWindowOperator) Name() string { return op.name }

// Processes one event at a time, but emits only when window is full
func (op *TumblingWindowOperator) Process(event model.Event) []model.Event {
    op.buffer = append(op.buffer, event)
    if len(op.buffer) >= op.size {
        out := op.inner.(interface {
            ProcessBatch([]model.Event) []model.Event
        }).ProcessBatch(op.buffer)
        op.buffer = nil
        return out
    }
    return nil
}
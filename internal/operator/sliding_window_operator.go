package operator

import (
    "goxstream/internal/model"
)

type SlidingWindowOperator struct {
    name      string
    size      int
    step      int
    buffer    []model.Event
    eventSeen int
    windowID  int
    inner     Operator // Should support ProcessBatch([]model.Event) []model.Event
}

func NewSlidingWindowOperator(name string, size, step int, inner Operator) *SlidingWindowOperator {
    return &SlidingWindowOperator{
        name:   name,
        size:   size,
        step:   step,
        buffer: []model.Event{},
        inner:  inner,
    }
}

func (op *SlidingWindowOperator) Name() string { return op.name }

func (op *SlidingWindowOperator) Process(event model.Event) []model.Event {
    op.buffer = append(op.buffer, event)
    op.eventSeen++
    if len(op.buffer) > op.size {
        op.buffer = op.buffer[1:]
    }

    if len(op.buffer) == op.size && ((op.eventSeen-op.size)%op.step == 0) {
        if batcher, ok := op.inner.(interface {
            ProcessBatch([]model.Event) []model.Event
        }); ok {
            op.windowID++
            out := batcher.ProcessBatch(op.buffer)
            // Annotate each result with the window ID
            for i := range out {
                if out[i].Data == nil {
                    out[i].Data = make(map[string]interface{})
                }
                out[i].Data["window_id"] = op.windowID
            }
            return out
        }
    }
    return nil
}

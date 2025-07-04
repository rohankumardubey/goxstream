package operator

import (
    "goxstream/internal/model"
    "time"
)

type TimeWindowWithWatermarkOperator struct {
    name            string
    windowDur       time.Duration
    allowedLateness time.Duration
    windows         map[time.Time][]model.Event
    maxEventTime    time.Time
    watermark       time.Time
    inner           Operator
    emittedWindows  map[time.Time]bool
    windowID        int
}

func NewTimeWindowWithWatermarkOperator(name string, windowDur, allowedLateness time.Duration, inner Operator) *TimeWindowWithWatermarkOperator {
    return &TimeWindowWithWatermarkOperator{
        name:            name,
        windowDur:       windowDur,
        allowedLateness: allowedLateness,
        windows:         make(map[time.Time][]model.Event),
        inner:           inner,
        emittedWindows:  make(map[time.Time]bool),
    }
}

func (op *TimeWindowWithWatermarkOperator) Name() string { return op.name }

func (op *TimeWindowWithWatermarkOperator) Process(event model.Event) []model.Event {
    ts := event.Timestamp
    if ts.After(op.maxEventTime) {
        op.maxEventTime = ts
        op.watermark = op.maxEventTime.Add(-op.allowedLateness)
    }
    // Assign event to its window
    windowEnd := ts.Truncate(op.windowDur).Add(op.windowDur)
    op.windows[windowEnd] = append(op.windows[windowEnd], event)

    out := []model.Event{}
    // Emit any windows whose end <= watermark (and not yet emitted)
    for end := range op.windows {
        if (end.Before(op.watermark) || end.Equal(op.watermark)) && !op.emittedWindows[end] {
            op.windowID++
            results := op.emitWindow(end)
            // annotate
            for i := range results {
                if results[i].Data == nil {
                    results[i].Data = make(map[string]interface{})
                }
                results[i].Data["window_end"] = end.Format(time.RFC3339)
                results[i].Data["window_id"] = op.windowID
            }
            out = append(out, results...)
            op.emittedWindows[end] = true
        }
    }
    return out
}

func (op *TimeWindowWithWatermarkOperator) emitWindow(windowEnd time.Time) []model.Event {
    events := op.windows[windowEnd]
    if batcher, ok := op.inner.(interface {
        ProcessBatch([]model.Event) []model.Event
    }); ok && len(events) > 0 {
        return batcher.ProcessBatch(events)
    }
    return nil
}

func (op *TimeWindowWithWatermarkOperator) Flush() []model.Event {
    out := []model.Event{}
    for end := range op.windows {
        if !op.emittedWindows[end] {
            op.windowID++
            results := op.emitWindow(end)
            for i := range results {
                if results[i].Data == nil {
                    results[i].Data = make(map[string]interface{})
                }
                results[i].Data["window_end"] = end.Format(time.RFC3339)
                results[i].Data["window_id"] = op.windowID
                results[i].Data["emitted_via_flush"] = true
            }
            out = append(out, results...)
            op.emittedWindows[end] = true
        }
    }
    return out
}


// -------------------- Basic Time-based Tumbling Window (no watermark) --------------------

type TimeWindowOperator struct {
    name      string
    windowDur time.Duration
    windowEnd time.Time
    buffer    []model.Event
    inner     Operator
    windowID  int
}

func NewTimeWindowOperator(name string, windowDur time.Duration, inner Operator) *TimeWindowOperator {
    return &TimeWindowOperator{
        name:      name,
        windowDur: windowDur,
        buffer:    []model.Event{},
        inner:     inner,
    }
}

func (op *TimeWindowOperator) Name() string { return op.name }

func (op *TimeWindowOperator) Process(event model.Event) []model.Event {
    if op.windowEnd.IsZero() {
        op.windowEnd = event.Timestamp.Truncate(op.windowDur).Add(op.windowDur)
    }

    out := []model.Event{}
    // If event is after windowEnd, emit and start new window(s)
    for !event.Timestamp.Before(op.windowEnd) {
        out = append(out, op.emitWindow()...)
        op.windowEnd = op.windowEnd.Add(op.windowDur)
    }
    op.buffer = append(op.buffer, event)
    return out
}

func (op *TimeWindowOperator) emitWindow() []model.Event {
    out := []model.Event{}
    if batcher, ok := op.inner.(interface {
        ProcessBatch([]model.Event) []model.Event
    }); ok && len(op.buffer) > 0 {
        op.windowID++
        result := batcher.ProcessBatch(op.buffer)
        for i := range result {
            if result[i].Data == nil {
                result[i].Data = make(map[string]interface{})
            }
            result[i].Data["window_end"] = op.windowEnd.Format(time.RFC3339)
            result[i].Data["window_id"] = op.windowID
        }
        out = append(out, result...)
    }
    op.buffer = nil
    return out
}

func (op *TimeWindowOperator) Flush() []model.Event {
    if len(op.buffer) == 0 {
        return nil
    }
    return op.emitWindow()
}

package operator

import (
    "goxstream/internal/model"
    "time"
)

type TimeSlidingWindowOperator struct {
    name       string
    windowSize time.Duration
    slide      time.Duration
    nextWindowEnd time.Time
    events     []model.Event
    inner      Operator // Should support ProcessBatch
    windowID   int
}

func NewTimeSlidingWindowOperator(name string, windowSize, slide time.Duration, inner Operator) *TimeSlidingWindowOperator {
    return &TimeSlidingWindowOperator{
        name:       name,
        windowSize: windowSize,
        slide:      slide,
        events:     []model.Event{},
        inner:      inner,
    }
}

func (op *TimeSlidingWindowOperator) Name() string { return op.name }

func (op *TimeSlidingWindowOperator) Process(event model.Event) []model.Event {
    op.events = append(op.events, event)
    out := []model.Event{}

    // On first event, initialize nextWindowEnd
    if op.nextWindowEnd.IsZero() {
        op.nextWindowEnd = event.Timestamp.Truncate(op.slide).Add(op.slide)
    }

    // While event.Timestamp >= nextWindowEnd, emit window and advance
    for !event.Timestamp.Before(op.nextWindowEnd) {
        // Emit window for [windowEnd - windowSize, windowEnd)
        windowStart := op.nextWindowEnd.Add(-op.windowSize)
        windowEvents := op.eventsInWindow(windowStart, op.nextWindowEnd)
        if len(windowEvents) > 0 {
            if batcher, ok := op.inner.(interface {
                ProcessBatch([]model.Event) []model.Event
            }); ok {
                op.windowID++
                windowResults := batcher.ProcessBatch(windowEvents)
                for i := range windowResults {
                    if windowResults[i].Data == nil {
                        windowResults[i].Data = make(map[string]interface{})
                    }
                    windowResults[i].Data["window_end"] = op.nextWindowEnd.Format(time.RFC3339)
                    windowResults[i].Data["window_id"] = op.windowID
                }
                out = append(out, windowResults...)
            }
        }
        op.nextWindowEnd = op.nextWindowEnd.Add(op.slide)
    }
    return out
}

// eventsInWindow returns events in [start, end)
func (op *TimeSlidingWindowOperator) eventsInWindow(start, end time.Time) []model.Event {
    res := []model.Event{}
    for _, e := range op.events {
        if !e.Timestamp.Before(start) && e.Timestamp.Before(end) {
            res = append(res, e)
        }
    }
    return res
}

// Call at end of input to flush remaining windows
func (op *TimeSlidingWindowOperator) Flush() []model.Event {
    out := []model.Event{}
    if len(op.events) == 0 || op.nextWindowEnd.IsZero() {
        return out
    }
    // Find the max timestamp to know where to stop
    maxT := op.events[0].Timestamp
    for _, e := range op.events {
        if e.Timestamp.After(maxT) {
            maxT = e.Timestamp
        }
    }
    for !op.nextWindowEnd.After(maxT.Add(op.windowSize)) {
        windowStart := op.nextWindowEnd.Add(-op.windowSize)
        windowEvents := op.eventsInWindow(windowStart, op.nextWindowEnd)
        if len(windowEvents) > 0 {
            if batcher, ok := op.inner.(interface {
                ProcessBatch([]model.Event) []model.Event
            }); ok {
                op.windowID++
                windowResults := batcher.ProcessBatch(windowEvents)
                for i := range windowResults {
                    if windowResults[i].Data == nil {
                        windowResults[i].Data = make(map[string]interface{})
                    }
                    windowResults[i].Data["window_end"] = op.nextWindowEnd.Format(time.RFC3339)
                    windowResults[i].Data["window_id"] = op.windowID
                }
                out = append(out, windowResults...)
            }
        }
        op.nextWindowEnd = op.nextWindowEnd.Add(op.slide)
    }
    return out
}

package operator

import (
	"fmt"
	"goxstream/internal/model"
)

type BatchReduceOperator struct {
	key string
	agg string // For now, only "count" is implemented
}

func NewBatchReduceOperator(key, agg string) *BatchReduceOperator {
	return &BatchReduceOperator{key: key, agg: agg}
}

func (op *BatchReduceOperator) Name() string { return "batch_reduce" }

// This method is used by window/batch processors.
func (op *BatchReduceOperator) ProcessBatch(events []model.Event) []model.Event {
	groups := make(map[string]int)
	for _, e := range events {
		k := fmt.Sprintf("%v", e.Data[op.key])
		groups[k]++
	}
	var result []model.Event
	for k, cnt := range groups {
		result = append(result, model.Event{Data: map[string]interface{}{
			op.key:  k,
			"count": cnt,
		}})
	}
	return result
}

// Process is not used in batch mode but needs to be defined.
func (op *BatchReduceOperator) Process(e model.Event) []model.Event {
	panic("BatchReduceOperator.Process should not be called; use ProcessBatch instead")
}

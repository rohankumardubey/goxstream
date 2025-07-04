package operator

import (
	"fmt"
	"goxstream/internal/model"
)

// --- Operator Factory Types ---

type OperatorFactory func(params map[string]interface{}) (Operator, error)

// --- Registry Map ---

var registry map[string]OperatorFactory

func init() {
	registry = map[string]OperatorFactory{
		"map":             mapOperatorFactory,
		"filter":          filterOperatorFactory,
		"reduce":          reduceOperatorFactory,
		"tumbling_window": tumblingWindowOperatorFactory,
		"sliding_window":  slidingWindowOperatorFactory,
	}
}

// --- Registry Entry Point ---

func BuildOperator(opSpec model.OperatorSpec) (Operator, error) {
	factory, ok := registry[opSpec.Type]
	if !ok {
		return nil, fmt.Errorf("unknown operator type: %s", opSpec.Type)
	}
	return factory(opSpec.Params)
}

// --- Helper: JSON numbers become float64, so cast to int ---

func toInt(val interface{}) (int, bool) {
	switch t := val.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	}
	return 0, false
}

// --- Map Operator Factory ---

func mapOperatorFactory(params map[string]interface{}) (Operator, error) {
	col, ok1 := params["col"].(string)
	val, ok2 := params["val"]
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("map operator expects col and val")
	}
	return NewMapOperator("map", func(e model.Event) model.Event {
		e.Data[col] = val
		return e
	}), nil
}

// --- Filter Operator Factory ---

func filterOperatorFactory(params map[string]interface{}) (Operator, error) {
	field, ok1 := params["field"].(string)
	eq, eqOk := params["eq"]
	if ok1 && eqOk {
		return NewFilterOperator("filter", func(e model.Event) bool {
			return e.Data[field] == eq
		}), nil
	}
	return nil, fmt.Errorf("filter operator expects field and eq")
}

// --- Reduce Operator Factory ---

func reduceOperatorFactory(params map[string]interface{}) (Operator, error) {
	key, ok1 := params["key"].(string)
	agg, ok2 := params["agg"].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("reduce operator expects key and agg")
	}
	return NewBatchReduceOperator(key, agg), nil
}

// --- Tumbling Window Operator Factory ---

func tumblingWindowOperatorFactory(params map[string]interface{}) (Operator, error) {
	size, ok1 := toInt(params["size"])
	innerSpec, ok2 := params["inner"].(map[string]interface{})
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("tumbling window expects size and inner")
	}
	innerType, ok3 := innerSpec["type"].(string)
	innerParams, ok4 := innerSpec["params"].(map[string]interface{})
	if !ok3 || !ok4 {
		return nil, fmt.Errorf("tumbling window inner must have type and params")
	}
	innerOp, err := BuildOperator(model.OperatorSpec{Type: innerType, Params: innerParams})
	if err != nil {
		return nil, fmt.Errorf("tumbling window inner op error: %w", err)
	}
	return NewTumblingWindowOperator("tumbling_window", size, innerOp), nil
}

// --- Sliding Window Operator Factory ---

func slidingWindowOperatorFactory(params map[string]interface{}) (Operator, error) {
	size, ok1 := toInt(params["size"])
	step, ok2 := toInt(params["step"])
	innerSpec, ok3 := params["inner"].(map[string]interface{})
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("sliding window expects size, step, inner operator")
	}
	innerType, ok4 := innerSpec["type"].(string)
	innerParams, ok5 := innerSpec["params"].(map[string]interface{})
	if !ok4 || !ok5 {
		return nil, fmt.Errorf("sliding window inner must have type and params")
	}
	innerOp, err := BuildOperator(model.OperatorSpec{Type: innerType, Params: innerParams})
	if err != nil {
		return nil, fmt.Errorf("sliding window inner op error: %w", err)
	}
	return NewSlidingWindowOperator("sliding_window", size, step, innerOp), nil
}

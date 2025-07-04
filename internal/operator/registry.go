package operator

import (
	"fmt"
	"goxstream/internal/model"
	"time"
)

type OperatorFactory func(params map[string]interface{}) (Operator, error)

var registry map[string]OperatorFactory

func init() {
	registry = map[string]OperatorFactory{
		"map":                   mapOperatorFactory,
		"filter":                filterOperatorFactory,
		"reduce":                reduceOperatorFactory,
		"tumbling_window":       tumblingWindowOperatorFactory,
		"sliding_window":        slidingWindowOperatorFactory,
		"time_window":           timeWindowOperatorFactory,           // basic time window (if you have it)
		"time_sliding_window":   timeSlidingWindowOperatorFactory,    // time-based sliding window
		"time_window_watermark": timeWindowWatermarkFactory,          // watermark support!
	}
}

func BuildOperator(opSpec model.OperatorSpec) (Operator, error) {
	factory, ok := registry[opSpec.Type]
	if !ok {
		return nil, fmt.Errorf("unknown operator type: %s", opSpec.Type)
	}
	return factory(opSpec.Params)
}

func toInt(val interface{}) (int, bool) {
	switch t := val.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	}
	return 0, false
}

// ----- Basic Operators -----

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

func reduceOperatorFactory(params map[string]interface{}) (Operator, error) {
	key, ok1 := params["key"].(string)
	agg, ok2 := params["agg"].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("reduce operator expects key and agg")
	}
	return NewBatchReduceOperator(key, agg), nil
}

// ----- Count-based Windows -----

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

// ----- Time-based Windows -----

func timeWindowOperatorFactory(params map[string]interface{}) (Operator, error) {
	durationStr, ok1 := params["duration"].(string)
	innerSpec, ok2 := params["inner"].(map[string]interface{})
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("time window expects duration and inner")
	}
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}
	innerType, ok3 := innerSpec["type"].(string)
	innerParams, ok4 := innerSpec["params"].(map[string]interface{})
	if !ok3 || !ok4 {
		return nil, fmt.Errorf("time window inner must have type and params")
	}
	innerOp, err := BuildOperator(model.OperatorSpec{Type: innerType, Params: innerParams})
	if err != nil {
		return nil, fmt.Errorf("time window inner op error: %w", err)
	}
	// This must match your operator's constructor for basic time window:
	return NewTimeWindowOperator("time_window", dur, innerOp), nil
}

func timeSlidingWindowOperatorFactory(params map[string]interface{}) (Operator, error) {
	windowSizeStr, ok1 := params["size"].(string)
	slideStr, ok2 := params["slide"].(string)
	innerSpec, ok3 := params["inner"].(map[string]interface{})
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("time sliding window expects size, slide, inner")
	}
	windowSize, err := time.ParseDuration(windowSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid window size: %w", err)
	}
	slide, err := time.ParseDuration(slideStr)
	if err != nil {
		return nil, fmt.Errorf("invalid slide: %w", err)
	}
	innerType, ok4 := innerSpec["type"].(string)
	innerParams, ok5 := innerSpec["params"].(map[string]interface{})
	if !ok4 || !ok5 {
		return nil, fmt.Errorf("time sliding window inner must have type and params")
	}
	innerOp, err := BuildOperator(model.OperatorSpec{Type: innerType, Params: innerParams})
	if err != nil {
		return nil, fmt.Errorf("time sliding window inner op error: %w", err)
	}
	return NewTimeSlidingWindowOperator("time_sliding_window", windowSize, slide, innerOp), nil
}

// ----- Watermarking Window -----

func timeWindowWatermarkFactory(params map[string]interface{}) (Operator, error) {
	durationStr, ok1 := params["duration"].(string)
	latenessStr, ok2 := params["allowed_lateness"].(string)
	innerSpec, ok3 := params["inner"].(map[string]interface{})
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("time window watermark expects duration, allowed_lateness, inner")
	}
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}
	lateness, err := time.ParseDuration(latenessStr)
	if err != nil {
		return nil, fmt.Errorf("invalid lateness: %w", err)
	}
	innerType, ok4 := innerSpec["type"].(string)
	innerParams, ok5 := innerSpec["params"].(map[string]interface{})
	if !ok4 || !ok5 {
		return nil, fmt.Errorf("inner must have type and params")
	}
	innerOp, err := BuildOperator(model.OperatorSpec{Type: innerType, Params: innerParams})
	if err != nil {
		return nil, fmt.Errorf("inner op error: %w", err)
	}
	return NewTimeWindowWithWatermarkOperator("time_window_watermark", dur, lateness, innerOp), nil
}

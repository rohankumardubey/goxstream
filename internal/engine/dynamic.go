package engine

import (
    "fmt"
    "goxstream/internal/model"
    "goxstream/internal/operator"
    "goxstream/internal/source"
    "goxstream/internal/sink"
)

func BuildAndRunPipeline(spec model.PipelineSpec) error {
    input := make(chan model.Event)
    output := make(chan model.Event)

    // Build operator chain
    var ops []operator.Operator
    for _, opSpec := range spec.Operators {
        op, err := operator.BuildOperator(opSpec)
        if err != nil {
            return fmt.Errorf("operator build error: %w", err)
        }
        ops = append(ops, op)
    }

    pipeline := Pipeline{Operators: ops}

    // Run pipeline in background
    go func() {
        pipeline.Run(input, output)
        close(output)
    }()

    // Sink
    go func() {
        sink.FileSink(spec.Sink.Path, output)
    }()

    // Source
    switch spec.Source.Type {
    case "file":
        return source.FileSource(spec.Source.Path, input)
    default:
        return fmt.Errorf("unsupported source type: %s", spec.Source.Type)
    }
}

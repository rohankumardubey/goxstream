package model

type PipelineSpec struct {
    Source    SourceSpec      `json:"source"`
    Operators []OperatorSpec  `json:"operators"`
    Sink      SinkSpec        `json:"sink"`
}

type SourceSpec struct {
    Type string `json:"type"` // e.g., "file"
    Path string `json:"path"`
	Raw map[string]interface{} `json:"-"`
}

type OperatorSpec struct {
    Type   string                 `json:"type"`   // e.g., "map", "filter"
    Params map[string]interface{} `json:"params"` // custom per operator
}

type SinkSpec struct {
    Type string `json:"type"` // e.g., "file"
    Path string `json:"path"`
	Raw map[string]interface{} `json:"-"`
}

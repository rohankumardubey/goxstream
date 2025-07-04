package model

import "time"

// Event represents a single record flowing through the pipeline.
type Event struct {
    Data      map[string]interface{}
    Timestamp time.Time
}

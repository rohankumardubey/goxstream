package operator

import "goxstream/internal/model"

type Operator interface {
    Name() string
    Process(event model.Event) []model.Event
}

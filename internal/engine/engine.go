package engine

import (
    "goxstream/internal/model"
    "goxstream/internal/operator"
)

type Pipeline struct {
    Operators []operator.Operator
}

func (p *Pipeline) Run(input <-chan model.Event, output chan<- model.Event) {
    for event := range input {
        events := []model.Event{event}
        for _, op := range p.Operators {
            next := []model.Event{}
            for _, e := range events {
                next = append(next, op.Process(e)...)
            }
            events = next
        }
        for _, out := range events {
            output <- out
        }
    }
}

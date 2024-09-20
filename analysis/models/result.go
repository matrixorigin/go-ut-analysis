package models

import (
	"time"
)

type Results []*Result
type Result struct {
	Name    string
	StartAt time.Time
	EndAt   time.Time
	Result  string
	Outputs []string
	Elapsed float64
}

func NewResult(e Event) *Result {
	return &Result{
		Name:    e.Test,
		StartAt: e.Time,
		Outputs: make([]string, 0, 100),
	}
}

func (t *Result) Finish(e Event) {
	t.Result = e.Action
	t.EndAt = e.Time
	t.Elapsed = e.Elapsed
}

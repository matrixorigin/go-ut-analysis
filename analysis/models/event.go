package models

import "time"

const (
	ActionStart  = "start"
	ActionRun    = "run"
	ActionPause  = "pause"
	ActionCont   = "cont"
	ActionPass   = "pass"
	ActionBench  = "bench"
	ActionFail   = "fail"
	ActionOutput = "output"
	ActionSkip   = "skip"
)

type Event struct {
	Time    time.Time `json:"Time"` // encodes as an RFC3339-format string
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"` // seconds
	Output  string    `json:"Output"`
}

func (te *Event) Reset() {
	te.Time = time.Time{}
	te.Package = ""
	te.Action = ""
	te.Test = ""
	te.Elapsed = 0
	te.Output = ""
}

func (te *Event) Copy() Event {
	return Event{
		Time:    te.Time,
		Action:  te.Action,
		Package: te.Package,
		Test:    te.Test,
		Elapsed: te.Elapsed,
		Output:  te.Output,
	}
}

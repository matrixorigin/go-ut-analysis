package test

import (
	"sort"
	"strings"
	"time"

	"github.com/matrixorigin/go-ut-analysis/analysis/models"
)

type Package struct {
	Name    string
	StartAt time.Time
	EndAt   time.Time
	Result  string
	Elapsed float64

	tests   map[string]*models.Result
	records models.Results

	succeed models.Results
	failed  models.Results
	skipped models.Results
}

func NewPackage(e models.Event) *Package {
	return &Package{
		Name:    e.Package,
		StartAt: e.Time,
		EndAt:   e.Time,
		Result:  models.ActionFail,
		tests:   make(map[string]*models.Result, 10),
		records: make(models.Results, 0, 10),
		succeed: make(models.Results, 0, 10),
		failed:  make(models.Results, 0, 10),
		skipped: make(models.Results, 0, 10),
	}
}

func (p *Package) ProcessEvent(e models.Event) {
	switch e.Action {
	case models.ActionRun:
		t := models.NewResult(e)
		p.records = append(p.records, t)
		p.tests[e.Test] = t
	case models.ActionOutput:
		p.addResultOutput(e)
	case models.ActionSkip, models.ActionPass, models.ActionFail:
		if e.Test != "" {
			p.addResultOutput(e)
			p.tests[e.Test].Finish(e)
			break
		}
		p.finish(e)
	}
}

func (p *Package) First(num int) models.Results {
	sort.Slice(p.records, func(i, j int) bool {
		return p.records[i].Elapsed > p.records[j].Elapsed
	})
	switch {
	case num <= 0 || num > len(p.records):
		return p.records
	default:
		return p.records[:num]
	}
}

func (p *Package) GetName() string {
	return p.Name
}

func (p *Package) GetElapsed() float64 {
	return p.Elapsed
}

func (p *Package) GetStartTime() time.Time {
	return p.StartAt.In(time.UTC)
}

func (p *Package) GetEndTime() time.Time {
	return p.EndAt.In(time.UTC)
}

func (p *Package) GetResult() string {
	return p.Result
}

func (p *Package) Get(name string) *models.Result {
	return p.tests[name]
}

func (p *Package) Success() models.Results {
	return p.succeed
}

func (p *Package) Fail() models.Results {
	return p.failed
}

func (p *Package) Skipped() models.Results {
	return p.skipped
}

func (p *Package) Finish() error {
	for _, t := range p.records {
		p.splitPush(t)
	}
	return nil
}

func (p *Package) splitPush(t *models.Result) {
	switch t.Result {
	case models.ActionFail:
		p.failed = append(p.failed, t)
	case models.ActionSkip:
		p.skipped = append(p.skipped, t)
	case models.ActionPass:
		p.succeed = append(p.succeed, t)
	}
}

func (p *Package) addResultOutput(e models.Event) {
	t := e.Test
	if ind := strings.LastIndexByte(t, '/'); ind != -1 {
		pe := e.Copy()
		pe.Test = t[:ind]
		p.addResultOutput(pe)
	}
	if _, ok := p.tests[e.Test]; !ok {
		p.tests[e.Test] = models.NewResult(e)
	}
	p.tests[e.Test].Outputs = append(p.tests[e.Test].Outputs, e.Output)
	p.tests[e.Test].EndAt = e.Time
	p.tests[e.Test].Elapsed = e.Time.Sub(p.tests[e.Test].StartAt).Seconds()
}

func (p *Package) finish(e models.Event) {
	p.EndAt = e.Time
	p.Result = e.Action
	p.Elapsed = e.Elapsed
}

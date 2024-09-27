package test

import (
	"sort"

	"github.com/matrixorigin/go-ut-analysis/analysis"
	"github.com/matrixorigin/go-ut-analysis/analysis/models"
)

type Packages struct {
	packages map[string]*Package
	records  []*Package

	success []analysis.Packager
	fail    []analysis.Packager
	skipped []analysis.Packager
}

func NewEmptyPackages() analysis.Analyzer {
	return &Packages{
		packages: make(map[string]*Package, 20),
		records:  make([]*Package, 0, 20),

		success: make([]analysis.Packager, 0, 20),
		fail:    make([]analysis.Packager, 0, 20),
		skipped: make([]analysis.Packager, 0, 20),
	}
}

func (ps *Packages) ProcessEvent(e models.Event) {
	switch e.Action {
	case models.ActionStart:
		p := NewPackage(e)
		ps.records = append(ps.records, p)
		ps.packages[e.Package] = p
	case models.ActionFail, models.ActionPass, models.ActionSkip, models.ActionRun, models.ActionOutput:
		ps.packages[e.Package].ProcessEvent(e)
	default:
		panic("unsupported action")
	}
}

func (ps *Packages) First(num int) []analysis.Packager {
	sort.Slice(ps.records, func(i, j int) bool {
		return ps.records[i].GetElapsed() > ps.records[j].GetElapsed()
	})
	length := num
	if num <= 0 || num > len(ps.records) {
		length = len(ps.records)
	}
	result := make([]analysis.Packager, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, ps.records[i])
	}
	return result
}

func (ps *Packages) Get(name string) analysis.Packager {
	return ps.packages[name]
}

func (ps *Packages) Success() []analysis.Packager {
	return ps.success
}

func (ps *Packages) Fail() []analysis.Packager {
	return ps.fail
}

func (ps *Packages) Skipped() []analysis.Packager {
	return ps.skipped
}

func (ps *Packages) Finish() error {
	for _, record := range ps.records {
		err := record.Finish()
		if err != nil {
			return err
		}
		ps.splitPush(record)
	}
	return nil

}

func (ps *Packages) splitPush(p *Package) {
	switch p.Result {
	case models.ActionPass:
		ps.success = append(ps.success, p)
	case models.ActionSkip:
		ps.skipped = append(ps.skipped, p)
	case models.ActionFail:
		ps.fail = append(ps.fail, p)
	}
}

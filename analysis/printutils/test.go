package printutils

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/matrixorigin/go-ut-analysis/analysis"
	"github.com/matrixorigin/go-ut-analysis/analysis/models"
	"github.com/matrixorigin/go-ut-analysis/util"
)

func NewPackageTable(title string, options ...Option) table.Writer {
	ops := GenDefaultOptions()
	for _, o := range options {
		o(&ops)
	}
	t := table.NewWriter()
	t.SetTitle(title)
	t.SetOutputMirror(io.MultiWriter(ops.writers...))
	if ops.packageSortBy != "" {
		t.SortBy([]table.SortBy{{Name: ops.packageSortBy}})
	}
	t.AppendHeader(DefaultPackageTile)
	return t
}

func WritePackageLine(t table.Writer, pkg analysis.Packager) {
	t.AppendRow(table.Row{
		pkg.GetName(),
		pkg.GetStartTime().Format(time.RFC3339),
		pkg.GetEndTime().Format(time.RFC3339),
		pkg.GetResult(),
		fmt.Sprintf("%.2f", pkg.GetElapsed()),
	})
	t.AppendSeparator()
}

func WriteTestArrayOutputs(ts models.Results, path string) error {
	for _, t := range ts {
		err := WriteTestOutputs(t, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteTestOutputs(tt *models.Result, path string) error {
	dst := util.MustCreateFile(path, util.ValidatedName(tt.Name))
	prefixTpl := "[TEST]: %s\n[START_AT]: %s\n[END_AT]: %s\n[RESULT]: %s\n[ELAPSED]: %.2f\n[OUTPUTS]:\n"
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf(prefixTpl, tt.Name, formatTime(tt.StartAt), formatTime(tt.EndAt), tt.Result, tt.Elapsed))
	for _, output := range tt.Outputs {
		builder.WriteString(output)
	}
	_, err := dst.Write([]byte(builder.String()))
	return err

}

func formatTime(t time.Time) string {
	return t.In(time.UTC).Format(time.RFC3339)
}

func NewTestTable(title string, options ...Option) table.Writer {
	ops := GenDefaultOptions()
	for _, o := range options {
		o(&ops)
	}
	t := table.NewWriter()
	t.SetTitle(title)
	t.SetOutputMirror(io.MultiWriter(ops.writers...))
	if ops.testSortBy != "" {
		t.SortBy([]table.SortBy{{Name: ops.testSortBy}})
	}
	t.AppendHeader(DefaultTestTile)
	return t
}

func WriteTestLine(t table.Writer, tt *models.Result) {
	t.AppendRow(table.Row{
		tt.Name,
		tt.StartAt.In(time.UTC).Format(time.RFC3339),
		tt.EndAt.In(time.UTC).Format(time.RFC3339),
		tt.Result,
		fmt.Sprintf("%.2f", tt.Elapsed),
	})
	t.AppendSeparator()
}

func PrintPackages(packages []analysis.Packager, title string, options ...Option) {
	t := NewPackageTable(title, options...)
	for _, pkg := range packages {
		WritePackageLine(t, pkg)
	}
	t.Render()
}

func PrintPackagesWithTests(packages []analysis.Packager, title string, options ...Option) {
	t := NewPackageTable(title, options...)
	tests := make([]table.Writer, 0, 10)
	for _, pkg := range packages {
		WritePackageLine(t, pkg)
		tt := NewTestTable(pkg.GetName(), options...)
		for _, result := range pkg.First(-1) {
			WriteTestLine(tt, result)
		}
		tests = append(tests, tt)
	}
	t.Render()
	for _, tt := range tests {
		tt.Render()
	}
}

func PrintTests(tests []*models.Result, title string, options ...Option) {
	tt := NewTestTable(title, options...)
	for _, result := range tests {
		WriteTestLine(tt, result)
	}
	tt.Render()
}

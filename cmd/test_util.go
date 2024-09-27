package cmd

import (
	"io"
	"path/filepath"
	"slices"
	"sync"

	"github.com/matrixorigin/go-ut-analysis/analysis"
	"github.com/matrixorigin/go-ut-analysis/analysis/printutils"
	"github.com/matrixorigin/go-ut-analysis/util"
)

func reportFailed(analyzer analysis.Analyzer, bases []io.Writer, reportPath string) {
	failed := analyzer.Fail()
	failedResultPath := filepath.Join(reportPath, "failed")
	mirrors := slices.Clone(bases)
	if len(reportPath) != 0 {
		f := util.MustCreateFile(failedResultPath, "result.txt")
		defer f.Close()
		mirrors = append(mirrors, f)
	}
	printutils.PrintPackages(failed, "Failed UT Packages", printutils.WithWriters(mirrors...))

	failedTestOutputsPath := filepath.Join(failedResultPath, "outputs")
	wait := sync.WaitGroup{}
	for _, p := range failed {
		wait.Add(1)
		go func() {
			defer wait.Done()
			err := printutils.WriteTestArrayOutputs(
				p.Fail(),
				filepath.Join(failedTestOutputsPath, filepath.Base(p.GetName())),
			)
			if err != nil {
				panic(err)
			}
		}()
	}
	wait.Wait()
}

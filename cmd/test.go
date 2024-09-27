package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/matrixorigin/go-ut-analysis/analysis/models"
	"github.com/matrixorigin/go-ut-analysis/analysis/printutils"
	"github.com/matrixorigin/go-ut-analysis/analysis/test"
	"github.com/matrixorigin/go-ut-analysis/util"
)

func initTestCommand() *cobra.Command {
	testCmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := util.MustGetStringFlag(cmd, "file")
			first := util.MustGetIntFlag(cmd, "first")
			showSuccess := util.MustGetBoolFlag(cmd, "success")
			showFailed := util.MustGetBoolFlag(cmd, "failed")
			showSkipped := util.MustGetBoolFlag(cmd, "skipped")
			showNoTest := util.MustGetBoolFlag(cmd, "no-test")

			reportPath := util.MustGetStringFlag(cmd, "report-path")
			stdout := util.MustGetBoolFlag(cmd, "stdout")

			bases := make([]io.Writer, 0, 2)
			if stdout {
				bases = append(bases, os.Stdout)
			}

			f := util.MustOpenFile(filePath)
			defer f.Close()
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			analyzer := test.NewEmptyPackages()
			event := &models.Event{}

			for scanner.Scan() {
				line := scanner.Bytes()
				event.Reset()
				err := json.Unmarshal(line, event)
				if err != nil {
					panic(err)
				}
				analyzer.ProcessEvent(*event)
			}
			err := analyzer.Finish()
			if err != nil {
				return err
			}
			// print top <first> Time-consuming UT Tests
			if first > 0 {
				mirrors := slices.Clone(bases)
				if len(reportPath) != 0 {
					topF := util.MustCreateFile(reportPath, "top.txt")
					mirrors = append(mirrors, topF)
				}
				printutils.PrintPackagesWithTests(analyzer.First(first), "The Most Time-consuming UT Tests",
					printutils.WithWriters(mirrors...),
				)
			}
			if showNoTest {
				skipped := analyzer.Skipped()
				mirrors := slices.Clone(bases)
				if len(reportPath) != 0 {
					skipF := util.MustCreateFile(reportPath, "no-test.txt")
					mirrors = append(mirrors, skipF)
				}
				printutils.PrintPackages(
					skipped,
					"No UT Test Packages",
					printutils.WithWriters(mirrors...))
			}
			if showSkipped {
				skipped := make(models.Results, 0, 100)
				all := analyzer.First(-1)
				for _, result := range all {
					skipped = append(skipped, result.Skipped()...)
				}
				mirrors := slices.Clone(bases)
				if len(reportPath) != 0 {
					topF := util.MustCreateFile(reportPath, "skipped.txt")
					mirrors = append(mirrors, topF)
				}
				printutils.PrintTests(skipped, "Skipped Tests", printutils.WithWriters(mirrors...))
			}
			if showSuccess {
				success := analyzer.Success()
				mirrors := slices.Clone(bases)
				if len(reportPath) != 0 {
					topF := util.MustCreateFile(reportPath, "success.txt")
					mirrors = append(mirrors, topF)
				}
				printutils.PrintPackages(success, "Successful UT Packages", printutils.WithWriters(mirrors...))
			}
			if showFailed {
				reportFailed(analyzer, bases, reportPath)
			}

			return nil
		},
	}

	testCmd.Flags().StringP("file", "f", "", "file path")
	testCmd.Flags().Int("first", 0, "show the top few tests that consume the most time")
	testCmd.Flags().String("report-path", "", "report path")
	testCmd.Flags().Bool("stdout", true, "write report to stdout")
	testCmd.Flags().Bool("success", true, "show success ut cases")
	testCmd.Flags().Bool("failed", true, "show failed ut cases")
	testCmd.Flags().Bool("skipped", true, "show skipped ut cases")
	testCmd.Flags().Bool("no-test", true, "show test files packages")

	return testCmd
}

package printutils

import "github.com/jedib0t/go-pretty/v6/table"

const (
	ColumnPackage   = "PACKAGE"
	ColumnTest      = "TEST"
	ColumnStartTime = "START_TIME"
	ColumnEndTime   = "END_TIME"
	ColumnResult    = "RESULT"
	ColumnElapsed   = "ELAPSED"
)

var (
	DefaultPackageTile = table.Row{ColumnPackage, ColumnStartTime, ColumnEndTime, ColumnResult, ColumnElapsed}
	DefaultTestTile    = table.Row{ColumnTest, ColumnStartTime, ColumnEndTime, ColumnResult, ColumnElapsed}
)

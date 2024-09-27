package analysis

import (
	"time"

	"github.com/matrixorigin/go-ut-analysis/analysis/models"
)

type Analyzer interface {
	ProcessEvent(e models.Event)
	First(num int) []Packager

	Get(name string) Packager
	Success() []Packager
	Fail() []Packager
	Skipped() []Packager

	Finish() error
}

type Packager interface {
	ProcessEvent(e models.Event)
	First(num int) models.Results
	GetElapsed() float64
	GetName() string
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetResult() string

	Get(name string) *models.Result
	Success() models.Results
	Fail() models.Results
	Skipped() models.Results

	Finish() error
}

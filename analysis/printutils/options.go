package printutils

import (
	"io"
	"os"
)

type Options struct {
	writers       []io.Writer
	packageSortBy string
	testSortBy    string
}

func GenDefaultOptions() Options {
	return Options{
		writers: []io.Writer{os.Stdout},
	}
}

type Option func(o *Options)

func WithPackageSortBy(sortBy string) Option {
	return func(o *Options) {
		if sortBy != "" {
			o.packageSortBy = sortBy
		}
	}
}

func WithTestSortBy(sortBy string) Option {
	return func(o *Options) {
		if sortBy != "" {
			o.testSortBy = sortBy
		}
	}
}

func WithWriters(writers ...io.Writer) Option {
	return func(o *Options) {
		if len(writers) != 0 {
			o.writers = writers
		}
	}
}

package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func MustGetFlag(cmd *cobra.Command, name string) *pflag.Flag {
	f := cmd.Flag(name)
	if f == nil {
		fmt.Printf("flag %s is not found", name)
		os.Exit(1)
	}
	return f
}

func MustGetStringFlag(cmd *cobra.Command, name string) string {
	return MustGetFlag(cmd, name).Value.String()
}

func MustGetBoolFlag(cmd *cobra.Command, name string) bool {
	f, err := cmd.Flags().GetBool(name)
	if err != nil {
		fmt.Printf("flag %s is not found", name)
		os.Exit(1)
	}
	return f
}

func MustGetIntFlag(cmd *cobra.Command, name string) int {
	f, err := cmd.Flags().GetInt(name)
	if err != nil {
		fmt.Printf("flag %s is not found", name)
		os.Exit(1)
	}
	return f
}

func MustCreatePath(path string) {
	if err := os.MkdirAll(path, fs.ModePerm); err != nil {
		panic(err)
	}
}

func ValidatedName(name string) string {
	dst := strings.Builder{}
	for _, i := range name {
		if ('A' <= i && i <= 'Z') || ('a' <= i && i <= 'z') || ('0' <= i && i <= '9') {
			dst.WriteRune(i)
			continue
		}
		dst.WriteRune('_')
	}
	return dst.String()
}

func MustCreateFile(path string, name string) *os.File {
	MustCreatePath(path)
	f, err := os.Create(filepath.Join(path, name))
	if err != nil {
		panic(err)
	}
	return f
}

func MustOpenFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

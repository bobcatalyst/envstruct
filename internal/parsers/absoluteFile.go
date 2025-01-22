package parsers

import "path/filepath"

type AbsFile struct{}

func (AbsFile) Name() string { return "absFile" }

func (AbsFile) Parse(s string) (string, error) {
    return filepath.Abs(s)
}

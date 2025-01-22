package parsers

import (
    "errors"
    "github.com/bobcatalyst/subflow"
    "regexp"
    "strings"
)

type Args struct{}

var argsPattern = regexp.MustCompile(`("[^"]*"|[^"\s]+)(\s+|$)`)

func (Args) Name() string { return "args" }

func (Args) Parse(s string) (subflow.CommandArgs, error) {
    var args []string
    for _, v := range argsPattern.FindAllString(s, -1) {
        if v = strings.TrimSpace(v); v != "" {
            args = append(args, v)
        }
    }

    if len(args) == 0 {
        return nil, errors.New("args array size must contain at least one element")
    }
    return subflow.NewCommandArgs(args[0], args[1:]), nil
}

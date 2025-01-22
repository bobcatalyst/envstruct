package parsers

import (
    "errors"
    "strconv"
)

type Port struct{}

func (Port) Name() string { return "port" }

func (Port) Parse(s string) (uint16, error) {
    i, err := strconv.ParseUint(s, 10, 16)
    if err != nil {
        return 0, err
    } else if i == 0 {
        return 0, errors.New("port cannot be 0")
    }
    return uint16(i), nil
}

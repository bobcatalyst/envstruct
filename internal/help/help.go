package help

import (
    "maps"
    "slices"
    "strings"
    "text/tabwriter"
)

type Help struct {
    envs map[string]*info
}

type info struct {
    Type    string
    Help    string
    Default *string
}

func (h *Help) Add(key, help, typeName string, def *string) {
    inf := info{
        Type:    typeName,
        Help:    help,
        Default: def,
    }

    if h.envs == nil {
        h.envs = map[string]*info{}
    }
    h.envs[key] = &inf
}

func (h *Help) String() string {
    var buf strings.Builder
    buf.WriteString("Environment:\n")
    tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

    for key := range slices.Values(slices.Sorted(maps.Keys(h.envs))) {
        inf := h.envs[key]
        typ := inf.Type
        if typ == "" {
            typ = "UNKNOWN"
        }
        line := []string{key + ":", "[" + typ + "]"}
        if inf.Default != nil {
            line = append(line, "(default: "+*inf.Default+")")
        } else {
            line = append(line, "")
        }
        line = append(line, inf.Help)
        tw.Write([]byte(strings.Join(line, "\t") + "\n"))
    }

    tw.Flush()
    return buf.String()
}

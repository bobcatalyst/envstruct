package help

import (
    "strings"
    "text/tabwriter"
)

type Help struct {
    envs []*info
}

type info struct {
    Key     string
    Type    string
    Help    string
    Default *string
}

func (h *Help) Add(key, help, typeName string, def *string) {
    inf := info{
        Key:     key,
        Type:    typeName,
        Help:    help,
        Default: def,
    }
    h.envs = append(h.envs, &inf)
}

func (h *Help) String() string {
    var buf strings.Builder
    buf.WriteString("Environment:\n")
    tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

    for _, inf := range h.envs {
        typ := inf.Type
        if typ == "" {
            typ = "UNKNOWN"
        }
        line := []string{inf.Key + ":", "[" + typ + "]"}
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

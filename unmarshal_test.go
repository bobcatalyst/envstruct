package envstruct

import (
    "testing"
)

func TestUnmarshal(t *testing.T) {
    type V struct {
        Goroot string        `env:"GOROOT" description:"go root" parser:"absFile"`
        Foo    string        `env:"FOO" description:"foo" default:"bar" parser:"string"`
        Help   func() string `env-help:""`
    }

    v, err := Unmarshal[V]()
    if err != nil {
        t.Errorf("error: %v", err)
    } else if v.Help() == "" {
        t.Errorf("no help")
    } else if v.Foo != "bar" {
        t.Errorf("foo not correct")
    } else if len(v.Goroot) == 0 {
        t.Errorf("goroot empty")
    }
    t.Logf("GOROOT: %s", v.Goroot)
    t.Logf("FOO: %s", v.Foo)
    t.Log("\n" + v.Help())
}

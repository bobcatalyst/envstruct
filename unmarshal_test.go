package envstruct

import (
    "math"
    "math/rand"
    "net"
    "testing"
)

func TestUnmarshal(t *testing.T) {
    type V struct {
        Foo  string        `env:"FOO" description:"foo" default:"bar" parser:"string"`
        Baz  string        `env:"BAZ" description:"baz" parser:"absFile"`
        Host net.IP        `env:"HOST" description:"host" parser:"ipv4"`
        Help func() string `env-help:""`
    }

    d := t.TempDir()
    t.Setenv("BAZ", d)

    r8 := func() uint8 { return uint8(rand.Intn(math.MaxUint8) + 1) }
    ip := net.IPv4(r8(), r8(), r8(), r8()).To4()
    t.Setenv("HOST", ip.String())

    v, err := Unmarshal[V]()
    if err != nil {
        t.Errorf("error: %v", err)
    } else if v.Help() == "" {
        t.Errorf("no help")
    } else if v.Foo != "bar" {
        t.Errorf("foo: not set to default")
    } else if v.Baz != d {
        t.Errorf("baz: path not parsed correctly")
    } else if !ip.Equal(v.Host) {
        t.Errorf("host: not parsed correctly")
    }
    t.Logf("FOO: %s", v.Foo)
    t.Logf("BAZ: %s", v.Baz)
    t.Logf("HOST: %s", v.Host)
    t.Log("\n" + v.Help())
}

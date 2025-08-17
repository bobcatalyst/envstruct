package envstruct

import (
    "errors"
    "github.com/bobcatalyst/envstruct/internal/parsers"
    "os"
    "reflect"
)

func init() {
    RegisterParser[parsers.String]()
    RegisterParser[parsers.AbsFile]()
    RegisterParser[parsers.Args]()
    RegisterParser[parsers.Port]()
    RegisterParser[parsers.IPv4Parser]()
}

type Parser[T any] interface {
    Name() string
    Parse(string) (T, error)
}

type Default[T any] interface {
    Parser[T]
    Default() string
}

var parserMap = map[string]func(dotEnv map[string]string, envKey string, defaultValue *string) (reflect.Value, error){}

var ErrNotFound = errors.New("env value not found")

func RegisterParser[P Parser[T], T any]() {
    parserMap[parserName[P, T]()] = parse[P, T]
}

func parse[P Parser[T], T any](dotEnv map[string]string, envKey string, defaultValue *string) (reflect.Value, error) {
    var v string
    var parser P
    if ev, ok := dotEnv[envKey]; ok {
        v = ev
    } else if ev, ok := os.LookupEnv(envKey); ok {
        v = ev
    } else if defaultValue != nil {
        v = *defaultValue
    } else if def, ok := any(parser).(Default[T]); ok {
        v = def.Default()
    } else {
        return reflect.Value{}, ErrNotFound
    }

    pv, err := parser.Parse(v)
    if err != nil {
        return reflect.Value{}, err
    }

    return reflect.ValueOf(pv), nil
}

func parserName[P Parser[T], T any]() string {
    var p P
    return p.Name()
}

package envstruct

import (
    "envstruct/internal/help"
    "envstruct/internal/parsers"
    "errors"
    "fmt"
    "github.com/joho/godotenv"
    "iter"
    "reflect"
)

const (
    TagEnv         = "env"
    TagEnvHelp     = "env-help"
    TagDescription = "description"
    TagDefault     = "default"
    TagParser      = "parser"
)

var (
    ErrInvalid          = errors.New("nil value")
    ErrNotStructPointer = errors.New("not a struct")
)

type TagOp string

type (
    OpReadTags struct{}
    TagOpParse struct{ Name string }
)

func (OpReadTags) String() string     { return "read tags" }
func (top TagOpParse) String() string { return fmt.Sprintf("parse env variable %q", top.Name) }

type ErrOpFail[Op fmt.Stringer] struct {
    Op    Op
    Field reflect.StructField
    Err   error
}

func (err *ErrOpFail[Op]) Unwrap() error {
    return err.Err
}

func (err *ErrOpFail[Op]) Error() string {
    return fmt.Sprintf("failed to %s for %q, %v", err.Op, err.Field.Name, err.Err)
}

type ErrTagNotSet struct {
    Name string
}

func (err *ErrTagNotSet) Error() string {
    return fmt.Sprintf("tag %q not set", err.Name)
}

type ErrParserNotFound struct {
    Name string
}

func (err *ErrParserNotFound) Error() string {
    return fmt.Sprintf("parser %q not found", err.Name)
}

type ErrMismatchedParseTypes struct {
    Parser reflect.Type
    Struct reflect.Type
}

func (err *ErrMismatchedParseTypes) Error() string {
    return fmt.Sprintf("cannot set struct field of type %q to type %q returned from parser", err.Struct.Name(), err.Parser.Name())
}

func Unmarshal[T any]() (v T, _ error) {
    dotEnv, _ := godotenv.Read()

    rv, err := getStruct[T]()
    if err != nil {
        return v, err
    }

    var helper help.Help
    var fieldParsers []func() error
    for st, sv := range iterStructValidFields(rv, TagEnv) {
        psr, err := readTag(dotEnv, st, sv, &helper)
        if err != nil {
            return v, err
        }
        fieldParsers = append(fieldParsers, psr)
    }

    setEnvHelp(rv, &helper)

    for _, parse := range fieldParsers {
        if err := parse(); err != nil {
            return v, err
        }
    }
    return rv.Interface().(T), nil
}

func readTag(dotEnv map[string]string, st reflect.StructField, sv reflect.Value, helper *help.Help) (func() error, error) {
    tags, err := unmarshalTags(st.Tag, sv.Kind())
    if err != nil {
        return nil, &ErrOpFail[OpReadTags]{
            Field: st,
            Err:   err,
        }
    }
    helper.Add(tags.key, tags.help, sv.Type().Name(), tags.def)

    return func() error {
        value, err := tags.Parse(dotEnv)
        if err != nil {
            return &ErrOpFail[TagOpParse]{
                Op: TagOpParse{
                    Name: tags.key,
                },
                Field: st,
                Err:   err,
            }
        } else if value.Type() != sv.Type() && !value.Type().Implements(sv.Type()) {
            return &ErrMismatchedParseTypes{
                Parser: value.Type(),
                Struct: sv.Type(),
            }
        }

        sv.Set(value)
        return nil
    }, nil
}

func setEnvHelp(rv reflect.Value, helper *help.Help) {
    for _, sv := range iterStructValidFields(rv, TagEnvHelp) {
        if sv.Kind() == reflect.Func {
            fn := sv.Type()
            if fn.NumIn() == 0 && fn.NumOut() == 1 && fn.Out(0) == reflect.TypeFor[string]() {
                sv.Set(reflect.ValueOf(helper.String))
                return
            }
        }
    }
}

func getStruct[T any]() (reflect.Value, error) {
    rv := reflect.New(reflect.TypeFor[T]()).Elem()
    if !rv.IsValid() {
        return reflect.Value{}, ErrInvalid
    } else if rv.Kind() != reflect.Struct {
        return reflect.Value{}, ErrNotStructPointer
    }
    return rv, nil
}

func iterStructValidFields(rv reflect.Value, tag string) iter.Seq2[reflect.StructField, reflect.Value] {
    return func(yield func(reflect.StructField, reflect.Value) bool) {
        for st, sv := range iterStruct(rv) {
            if fieldValid(st, sv, tag) {
                if !yield(st, sv) {
                    return
                }
            }
        }
    }
}

func iterStruct(rv reflect.Value) iter.Seq2[reflect.StructField, reflect.Value] {
    rt := rv.Type()
    return func(yield func(reflect.StructField, reflect.Value) bool) {
        for i := 0; i < rv.NumField(); i++ {
            if !yield(rt.Field(i), rv.Field(i)) {
                return
            }
        }
    }
}

func fieldValid(st reflect.StructField, sv reflect.Value, tag string) bool {
    _, ok := st.Tag.Lookup(tag)
    return st.IsExported() && sv.CanSet() && ok
}

type tagValues struct {
    key    string
    help   string
    parser string
    def    *string
}

func (t *tagValues) Parse(dotEnv map[string]string) (value reflect.Value, err error) {
    parser, ok := parserMap[t.parser]
    if !ok {
        return value, &ErrParserNotFound{Name: t.parser}
    }
    return parser(dotEnv, t.key, t.def)
}

func unmarshalTags(tags reflect.StructTag, kind reflect.Kind) (_ *tagValues, err error) {
    t := new(tagValues)
    var ok bool
    if t.parser, ok = tags.Lookup(TagParser); !ok {
        t.parser, err = parserForKind(kind)
    }

    if t.key, ok = tags.Lookup(TagEnv); !ok {
        err = errors.Join(err, &ErrTagNotSet{Name: TagEnv})
    }

    if t.help, ok = tags.Lookup(TagDescription); !ok {
        err = errors.Join(err, &ErrTagNotSet{Name: TagDescription})
    }

    if d, ok := tags.Lookup(TagDefault); ok {
        t.def = &d
    }
    return t, err
}

func parserForKind(kind reflect.Kind) (string, error) {
    var p interface{ Name() string }
    switch kind {
    case reflect.String:
        p = parsers.String{}
    default:
        return "", &ErrTagNotSet{Name: TagParser}
    }
    return p.Name(), nil
}

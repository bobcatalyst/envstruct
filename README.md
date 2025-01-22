# EnvStruct

EnvStruct simplifies configuring applications via environment variables. It provides tools to unmarshal structured configuration directly from the environment into a struct, with support for custom parsers and help documentation.

## Installation

To use EnvStruct, install it via `go get`:

```bash
go get github.com/bobcatalyst/envstruct
```

## Usage

### Defining a Configuration Struct

Define your application's configuration as a Go struct with `env` tags for environment variable names:

```go
package main

import (
    "envstruct"
    "fmt"
)

type Config struct {
    Goroot string        `env:"GOROOT" description:"Go root directory" parser:"absFile"`
    Foo    string        `env:"FOO" description:"An example variable" default:"bar" parser:"string"`
    Help   func() string `env-help:""`
}

func main() {
    cfg, err := envstruct.Unmarshal[Config]()
    if err != nil {
        panic(err)
    }

    fmt.Println("GOROOT:", cfg.Goroot)
    fmt.Println("FOO:", cfg.Foo)
    fmt.Println("Help:", cfg.Help())
}
```

### Supported Parsers

EnvStruct includes several built-in parsers:

- `string`: Parses strings (default parser).
- `absFile`: Resolves file paths to their absolute form.
- `args`: Parses command-line arguments.
- `ipv4`: Validates and parses IPv4 addresses.
- `port`: Validates and parses network port numbers.

### Adding Custom Parsers

You can add your own parsers by implementing the `Parser` interface:

```go
package parsers

type MyCustomParser struct{}

func (MyCustomParser) Name() string { return "myCustomParser" }
func (MyCustomParser) Parse(s string) (MyType, error) {
    // Custom parsing logic here
}
```

Then register the parser:

```go
package main

import (
    "envstruct"
    "yourmodule/parsers"
)

func init() {
    envstruct.RegisterParser[parsers.MyCustomParser]()
}
```

### Generating Help Documentation

EnvStruct automatically generates help documentation for environment variables. The `env-help` tag allows mapping a function to retrieve the help string:

```go
Help   func() string `env-help:""`
```

Call the function in your application to display usage information:

```go
fmt.Println(cfg.Help())
```

### Using `.env` Files

Place your environment variables in a `.env` file, and EnvStruct will load them automatically:

```
GOROOT=/usr/local/go
FOO=myValue
```

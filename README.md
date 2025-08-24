[![License](https://img.shields.io/badge/license-MIT-blue)](https://choosealicense.com/licenses/mit/)
![go.mod Go version](https://img.shields.io/github/go-mod/go-version/krostar/cli?label=go)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/krostar/cli)
[![Latest tag](https://img.shields.io/github/v/tag/krostar/cli)](https://github.com/krostar/cli/tags)
[![Go Report](https://goreportcard.com/badge/github.com/krostar/cli)](https://goreportcard.com/report/github.com/krostar/cli)

# cli: Command Line Interface Made Simpler

The `cli` package provides a framework-agnostic way to define command-line interfaces in Go.
Its primary goal is to decouple your application's CLI logic from specific CLI frameworks like `spf13/cobra` or `urfave/cli`.
This allows for easier testing, greater flexibility, and the ability to switch between underlying CLI frameworks without significant code changes.

## Motivation

Existing CLI frameworks, while powerful, often lead to tight coupling.
Changing frameworks can require extensive rewrites.

`cli` aims to solve this by:

- **Abstraction**: Providing a simple, consistent interface for defining commands.
- **Testability**: Making it easy to test command logic in isolation, without framework dependencies.
- **Flexibility**: Allowing you to choose (and change) the underlying CLI framework that best suits your needs.
- **Simplicity**: Focusing on core CLI functionality, avoiding unnecessary complexity.

## API Stability

This project follows [Semantic Versioning](https://semver.org/).

## Usage

### Basic Command

The simplest command implements the `Command` interface:

```go
package main

import (
    "context"
    "os"

    "github.com/krostar/cli"
    spf13cobra "github.com/krostar/cli/mapper/spf13/cobra"
)

type myCommand struct{}

func (myCommand) Execute(ctx context.Context, args []string, dashedArgs []string) error {
    // Your command logic here.
    // args: Positional arguments.
    // dashedArgs: Arguments after a "--".
    return nil
}

func main() {
    cmd := cli.New(myCommand{})
    err := spf13cobra.Execute(context.Background(), os.Args, cmd)
    cli.Exit(context.Background(), err)
}
```

### Adding Subcommands

```go
type subCommand struct{}

func (subCommand) Execute(ctx context.Context, args []string, dashedArgs []string) error {
    // Subcommand logic.
    return nil
}

func main() {
    cmd := cli.New(myCommand{}).AddCommand("sub", subCommand{}) // Add a subcommand named "sub".

    // Or, mount a complete CLI as a subcommand:
    subCLI := cli.New(subCommand{})
    cmd.Mount("another", subCLI)

    err := spf13cobra.Execute(context.Background(), os.Args, cmd)
    cli.Exit(context.Background(), err)
}
```

### Flags

```go
type flagCommand struct {
    name string
    age  int
    tags []string
}

func (c *flagCommand) Flags() []cli.Flag {
    return []cli.Flag{
        cli.NewBuiltinFlag("name", "", &c.name, "Your name"),
        cli.NewBuiltinFlag("age", "", &c.age, "Your age"),
        cli.NewBuiltinSliceFlag("tags", "t", &c.tags, "Comma-separated tags"),
    }
}

func (c *flagCommand) Execute(ctx context.Context, args []string, dashedArgs []string) error {
    // Access flag values: c.name, c.age, c.tags
    return nil
}
```

### Hooks

```go
type hookedCommand struct{}

func (hookedCommand) Execute(ctx context.Context, args []string, dashedArgs []string) error {
    // command logic
    return nil
}

func (c hookedCommand) Hook() *cli.Hook {
    return &cli.Hook{
        BeforeCommandExecution: func(ctx context.Context) error {
            // Code to run before Execute.
            return nil
        },
        AfterCommandExecution: func(ctx context.Context) error {
            // Code to run after Execute.
            return nil
        },
    }
}
```

### Signal Handling

```go
func main() {
    ctx, cancel := cli.NewContextCancelableBySignal(syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    err := spf13cobra.Execute(ctx, os.Args, /* ... */)
    cli.Exit(ctx, err)
}
```

## Configuration Management

The `cli` package provides a powerful configuration system through the `cfg` package that allows loading configuration from multiple sources with precedence.

### Configuration Sources

The following sources are supported out of the box:

- **Default Values**: Set default values for your configuration
- **Environment Variables**: Load configuration from environment variables
- **Configuration Files**: Load configuration from YAML/JSON files
- **Command-line Flags**: Load configuration from command-line flags

### Configuration Example

```go
import (
    "github.com/krostar/cli"
    clicfg "github.com/krostar/cli/cfg"
    "github.com/krostar/cli/cfg/source/default"
    "github.com/krostar/cli/cfg/source/env"
    "github.com/krostar/cli/cfg/source/file"
    "github.com/krostar/cli/cfg/source/flag"
)

// Config structure with environment variable mappings
type Config struct {
    Server struct {
        Host string `env:"SERVER_HOST"`
        Port int    `env:"SERVER_PORT"`
    }
    LogLevel   string `env:"LOG_LEVEL"`
    ConfigFile string `env:"CONFIG_FILE"`
}

// SetDefault implements the default values
func (cfg *Config) SetDefault() {
    cfg.Server.Host = "localhost"
    cfg.Server.Port = 8080
    cfg.LogLevel = "info"
    cfg.ConfigFile = "config.yaml"
}

// Command with configuration
type MyCommand struct {
    config Config
}

// Define flags that map to your configuration
func (cmd *MyCommand) Flags() []cli.Flag {
    return []cli.Flag{
        cli.NewBuiltinFlag("config", "c", &cmd.config.ConfigFile, "Path to config file"),
        cli.NewBuiltinFlag("host", "", &cmd.config.Server.Host, "Server host"),
        cli.NewBuiltinFlag("port", "p", &cmd.config.Server.Port, "Server port"),
        cli.NewBuiltinFlag("log-level", "l", &cmd.config.LogLevel, "Log level"),
    }
}

// Use hooks to load configuration in order of precedence
func (cmd *MyCommand) Hook() *cli.Hook {
    return &cli.Hook{
        BeforeCommandExecution: clicfg.BeforeCommandExecutionHook(
            &cmd.config,
            // Sources are applied in order, with later sources taking precedence
            sourcedefault.Source[Config](),                              // 1. Defaults
            sourcefile.Source(getConfigFilePath, yamlUnmarshaler, true), // 2. Config file
            sourceenv.Source[Config]("APP"),                             // 3. Environment variables
            sourceflag.Source[Config](cmd),                              // 4. Command-line flags
        ),
    }
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

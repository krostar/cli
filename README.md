[![License](https://img.shields.io/badge/license-MIT-blue)](https://choosealicense.com/licenses/mit/)
![go.mod Go version](https://img.shields.io/github/go-mod/go-version/krostar/cli?label=go)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/krostar/cli)
[![Latest tag](https://img.shields.io/github/v/tag/krostar/cli)](https://github.com/krostar/cli/tags)
[![Go Report](https://goreportcard.com/badge/github.com/krostar/cli)](https://goreportcard.com/report/github.com/krostar/cli)

# Command Line Interface made simpler

The main goal of this package is to avoid being too tightly coupled with existing CLI framework.

## Motivation

As of today, there are some very nice frameworks to use to handle command line interface; one of them is [spf13/cobra](https://github.com/spf13/cobra);
but it is not super obvious how to create something decoupled from cobra.

Why would someone want to do that ?
- i was using [spf13/cobra](https://github.com/spf13/cobra) and wanted to try [urfave/cli](https://github.com/urfave/cli)
    and it implied quite a lot of changes to my application
- i tried to avoid using a framework, but I was reinventing a lot of stuff
- frameworks expose you to a lot of features, but it comes with a lot of concepts, and two different frameworks does not necessarily come with the same concepts nor the same implementation

I did not want to recreate yet a new CLI framework to solve my problem because existing ones already are complete,
but I am not using a lot of features and I wanted to keep things simple to use, simple to test, and simple to extend.

## Usage

The simplest command is defined by implementing the `Command` interface
```go
type myCommand struct{}
func (myCommand) Execute(ctx context.Context, args []string, dashedArgs []string) error {
    return nil
}

func main() {
    cmd := cli.NewCommand(myCommand{})
    err := spf13cobra.Execute(context.Background(), os.Args, cmd)
    cli.Exit(err)
}
```

This will create a CLI with one root command. This CLI is then mapped to be executed by the spf13/cobra framework.

A more useful / complex example can be found in `internal/example`.

package cli

import (
	"context"
	"io"
	"os"
)

// CLI stores the settings associated to the CLI.
type CLI struct {
	Name        string
	Command     Command
	SubCommands []*CLI
}

// NewCommand creates a new CLI builder.
func NewCommand(name string, command Command) *CLI {
	return &CLI{
		Name:    name,
		Command: command,
	}
}

// AddCommand adds a new subcommand to the CLI.
func (cli *CLI) AddCommand(name string, cmd Command) *CLI {
	cli.SubCommands = append(cli.SubCommands, NewCommand(name, cmd))
	return cli
}

// Add adds a whole new CLI as a subcommand of the CLI.
func (cli *CLI) Add(sub *CLI) *CLI {
	cli.SubCommands = append(cli.SubCommands, sub)
	return cli
}

// Command defines the minimal interface required to execute a CLI command.
type Command interface {
	Execute(ctx context.Context, args []string, dashedArgs []string) error
}

// Hooks defines multiple entry point to add behavior to the CLI lifecycle.
type Hooks struct {
	BeforeFlagsDefinition            func(ctx context.Context) error
	BeforeCommandExecution           func(ctx context.Context) error
	AfterCommandExecution            func(ctx context.Context) error
	PersistentBeforeCommandExecution func(ctx context.Context) error
	PersistentAfterCommandExecution  func(ctx context.Context) error
}

// SetExitLogger sets the logger used by the CLI to write the exit message if any.
func SetExitLogger(ctx context.Context, writer io.WriteCloser) {
	SetMetadata(ctx, ctxExitLogger, writer)
}

func getExitLogger(ctx context.Context) io.WriteCloser {
	rawWriter := GetMetadata(ctx, ctxExitLogger)
	if writer, ok := rawWriter.(io.WriteCloser); ok {
		return writer
	}
	return os.Stderr
}

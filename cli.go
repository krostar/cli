package cli

import (
	"context"
)

type CLI struct {
	Name        string
	Command     Command
	SubCommands []*CLI
}

func NewCommand(name string, command Command) *CLI {
	return &CLI{
		Name:    name,
		Command: command,
	}
}

func (cli *CLI) AddCommand(name string, cmd Command) *CLI {
	cli.SubCommands = append(cli.SubCommands, NewCommand(name, cmd))
	return cli
}

func (cli *CLI) Add(sub *CLI) *CLI {
	cli.SubCommands = append(cli.SubCommands, sub)
	return cli
}

type Command interface {
	Execute(ctx context.Context, args []string, dashedArgs []string) error
}

type Hooks struct {
	BeforeFlagsDefinition            func(ctx context.Context) error
	BeforeCommandExecution           func(ctx context.Context) error
	AfterCommandExecution            func(ctx context.Context) error
	PersistentBeforeCommandExecution func(ctx context.Context) error
	PersistentAfterCommandExecution  func(ctx context.Context) error
}

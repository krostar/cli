package example

import (
	"context"
	"errors"
	"fmt"

	"github.com/krostar/cli"
)

type CommandPrint struct{}

func (cmd *CommandPrint) Usage() string {
	return "pos args -- dashed args"
}

func (cmd *CommandPrint) Examples() []string {
	return []string{
		"print a",
		"print -- b",
		"print a -- b",
	}
}

func (cmd *CommandPrint) Description() string {
	return "print positional and dashed arguments\n" +
		"print prints at least one and maximum three arguments, and a unlimited number of dashed arguments"
}

func (cmd CommandPrint) Execute(_ context.Context, args []string, dashedArgs []string) error {
	if len(args) == 0 {
		return cli.ErrorShowHelp(errors.New("there should be at least 1 arg to print"))
	}
	if len(args) > 3 {
		return errors.New("there should be no more than 3 args to print")
	}

	for i, arg := range args {
		fmt.Printf("args[%d] = %s", i, arg)
	}

	for i, arg := range dashedArgs {
		fmt.Printf("dashedArgs[%d] = %s", i, arg)
	}

	return nil
}

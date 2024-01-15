package spf13cobra

import (
	"context"

	"github.com/krostar/cli"
)

// Execute executes the CLI with the spf13/cobra backend. Provided args are supposed to contain application name.
func Execute(ctx context.Context, args []string, c *cli.CLI) error {
	command, err := buildCobraCommandFromCLIRecursively(ctx, c)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		args = args[1:]
	}

	command.SetArgs(args)
	return command.Execute()
}

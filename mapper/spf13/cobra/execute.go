package spf13cobra

import (
	"context"
	"fmt"

	"github.com/krostar/test"
	"github.com/spf13/cobra"

	"github.com/krostar/cli"
	mapper "github.com/krostar/cli/mapper/internal"
)

// Execute executes the CLI with the spf13/cobra backend.
// It builds a Cobra command tree from the provided cli.CLI instance,
// sets the command arguments, applies any provided options, and executes the command.
//
// Note: The first argument in args (if present) is used as the CLI name and
// removed from the argument list passed to the actual command.
func Execute(ctx context.Context, args []string, c *cli.CLI, opts ...Option) error {
	// set CLI name from the first argument (typically the binary name)
	// and remove it from the arguments passed to the command
	if c.Name == "" && len(args) > 0 {
		c.Name = args[0]
		args = args[1:]
	}

	command, err := buildCobraCommandFromCLIRecursively(ctx, c)
	if err != nil {
		return fmt.Errorf("unable not build cobra command from cli: %w", err)
	}

	command.SetArgs(args)

	for _, opt := range opts {
		opt(command)
	}

	return command.Execute()
}

// Option is a function type for configuring a cobra.Command before execution.
// This allows for customizing various aspects of the command behavior.
type Option func(p *cobra.Command)

// ForTest returns an Option that configures a cobra.Command for testing purposes.
// It redirects both standard output and error streams to the test's logging system
// with a "[CLI]: " prefix, making CLI output clearly identifiable in test logs.
func ForTest(t test.TestingT) Option {
	writer := mapper.TestWriter(t)
	return func(p *cobra.Command) {
		p.SetOut(writer)
		p.SetErr(writer)
	}
}

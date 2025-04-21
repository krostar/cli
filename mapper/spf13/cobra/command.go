package spf13cobra

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/krostar/cli"
	mapper "github.com/krostar/cli/mapper/internal"
)

// buildCobraCommandFromCLIRecursively constructs a `cobra.Command` from a `cli.CLI` instance.
// It recursively processes subcommands, creating a tree of `cobra.Command`s that mirrors the
// structure of the `cli.CLI`.
func buildCobraCommandFromCLIRecursively(ctx context.Context, c *cli.CLI) (*cobra.Command, error) {
	ctx = cli.NewCommandContext(ctx)
	ctx = mapper.Context(c.Command, ctx)

	command, err := buildCobraCommandFromCLICommand(ctx, c.Name, c.Command)
	if err != nil {
		return nil, fmt.Errorf("unable to build spf13/cobra command %s: %w", c.Name, err)
	}

	for _, subCommand := range c.SubCommands {
		subCommand, err := buildCobraCommandFromCLIRecursively(ctx, subCommand)
		if err != nil {
			return nil, fmt.Errorf("unable to build spf13/cobra sub-command %s of command %s: %w", subCommand.Name(), c.Name, err)
		}
		command.AddCommand(subCommand)
	}

	return command, nil
}

// buildCobraCommandFromCLICommand creates a single `cobra.Command` from a `cli.Command`.
func buildCobraCommandFromCLICommand(ctx context.Context, commandName string, cliCommand cli.Command) (*cobra.Command, error) {
	var commandExample string
	if examples := mapper.Examples(cliCommand); len(examples) > 0 {
		commandExample = "  " + strings.Join(mapper.Examples(cliCommand), "\n  ")
	}

	cobraCommand := &cobra.Command{
		Use:     commandName + " " + mapper.Usage(cliCommand),
		Short:   mapper.ShortDescription(cliCommand),
		Long:    mapper.Description(cliCommand),
		Example: commandExample,
		RunE:    cobraHandlerFromCLIHandler(ctx, cliCommand),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
		},
		SilenceErrors:         true,
		SilenceUsage:          true,
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
	}

	if err := setCobraHooksFromCLIHooks(ctx, cobraCommand, mapper.Hook(cliCommand), mapper.PersistentHook(cliCommand)); err != nil {
		return nil, err
	}

	localFlags, persistentFlags := mapper.Flags(cliCommand), mapper.PersistentFlags(cliCommand)
	cli.SetInitializedFlagsInContext(ctx, localFlags, persistentFlags)

	setCobraFlagsFromCLIFlags(cobraCommand.Flags(), localFlags)
	setCobraFlagsFromCLIFlags(cobraCommand.PersistentFlags(), persistentFlags)

	return cobraCommand, nil
}

// getCommandArguments separates the arguments passed to a command into positional arguments
// and dashed arguments (arguments after "--"). It uses the `ArgsLenAtDash` method of the
// `cobra.Command` to determine the split point.
func getCommandArguments(command *cobra.Command, args []string) ([]string, []string) {
	switch argsSeparatedAt := command.ArgsLenAtDash(); {
	case len(args) == 0:
		return nil, nil
	case argsSeparatedAt == 0 && len(args) > 0:
		return nil, args
	case argsSeparatedAt > 0 && len(args[argsSeparatedAt:]) > 0:
		return args[:argsSeparatedAt], args[argsSeparatedAt:]
	default:
		return args, nil
	}
}

// cobraHandlerFromCLIHandler adapts a `cli.Command`'s `Execute` method to the `cobra.Command`'s `RunE` function signature.
// It handles the argument splitting and calls the `Execute` method with the appropriate context and arguments.
// It also handles the `ShowHelpError`, displaying the command's usage if required.
func cobraHandlerFromCLIHandler(ctx context.Context, cmd cli.Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		args, dashedArgs := getCommandArguments(c, args)

		err := cmd.Execute(ctx, args, dashedArgs)

		var showHelpErr cli.ShowHelpError
		if errors.As(err, &showHelpErr) {
			if showHelpErr.ShowHelp() {
				err = multierr.Append(err, c.Usage())
			}
		}

		return err
	}
}

// setCobraHooksFromCLIHooks sets the pre-run and post-run hooks for a `cobra.Command`
// based on the `cli.Hook` and `cli.PersistentHook` provided. It ensures that persistent
// hooks are executed in the correct order (parent first, then child).
func setCobraHooksFromCLIHooks(ctx context.Context, c *cobra.Command, hook *cli.Hook, persistentHook *cli.PersistentHook) error {
	if err := persistentHook.BeforeFlagsDefinition(ctx); err != nil {
		return fmt.Errorf("pre-flag-definition hook failed: %w", err)
	}

	c.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PersistentPreRunE != nil {
			if err := parent.PersistentPreRunE(parent, args); err != nil {
				return err
			}
		}
		return persistentHook.BeforeCommandExecution(ctx)
	}

	c.PersistentPostRunE = func(c *cobra.Command, args []string) error {
		if err := persistentHook.AfterCommandExecution(ctx); err != nil {
			return err
		}

		if parent := c.Parent(); parent != nil && parent.PersistentPostRunE != nil {
			if err := parent.PersistentPostRunE(parent, args); err != nil {
				return err
			}
		}

		return nil
	}

	c.PreRunE = func(*cobra.Command, []string) error { return hook.BeforeCommandExecution(ctx) }
	c.PostRunE = func(*cobra.Command, []string) error { return hook.AfterCommandExecution(ctx) }

	return nil
}

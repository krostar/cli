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

func buildCobraCommandFromCLIRecursively(ctx context.Context, c *cli.CLI) (*cobra.Command, error) {
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

	setCobraFlagsFromCLIFlags(cobraCommand.Flags(), mapper.Flags(cliCommand))
	setCobraFlagsFromCLIFlags(cobraCommand.PersistentFlags(), mapper.PersistentFlags(cliCommand))

	return cobraCommand, nil
}

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
		if parent := c.Parent(); parent != nil && parent.PersistentPostRunE != nil {
			if err := parent.PersistentPostRunE(parent, args); err != nil {
				return err
			}
		}
		return persistentHook.AfterCommandExecution(ctx)
	}

	c.PreRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PreRunE != nil {
			if err := parent.PreRunE(parent, args); err != nil {
				return err
			}
		}
		return hook.BeforeCommandExecution(ctx)
	}

	c.PostRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PostRunE != nil {
			if err := parent.PostRunE(parent, args); err != nil {
				return err
			}
		}
		return hook.AfterCommandExecution(ctx)
	}

	return nil
}

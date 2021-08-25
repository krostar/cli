package cobra

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
	"github.com/krostar/cli/mapper"
)

// Execute executes the CLI with the cobra backend.
func Execute(ctx context.Context, c *cli.CLI, args []string, opts ...Option) error {
	options := newOptions()
	for _, o := range opts {
		o(options)
	}

	command, err := buildCommandRecursively(ctx, c, options)
	if err != nil {
		return err
	}
	command.Use = app.Name()
	command.Version = fmt.Sprintf("%s, compiled %s", app.Version(), app.BuiltAt().Local().Format(time.RFC3339))

	if len(args) > 0 {
		args = args[1:]
	}
	command.SetArgs(args)

	return command.Execute()
}

func buildCommandRecursively(ctx context.Context, cli *cli.CLI, options *options) (*cobra.Command, error) {
	command, err := buildCommand(ctx, cli.Name, cli.Command, options)
	if err != nil {
		return nil, fmt.Errorf("unable to build spf13/cobra command %s: %w", cli.Name, err)
	}

	for _, sub := range cli.SubCommands {
		subCommand, err := buildCommandRecursively(ctx, sub, options)
		if err != nil {
			return nil, fmt.Errorf("unable to build spf13/cobra sub-command %s of command %s: %w", sub.Name, cli.Name, err)
		}
		command.AddCommand(subCommand)
	}

	return command, nil
}

func buildCommand(ctx context.Context, commandName string, cmd cli.Command, options *options) (*cobra.Command, error) {
	var example string
	if examples := mapper.Examples(cmd); examples != nil {
		example = "  " + strings.Join(mapper.Examples(cmd), "\n  ")
	}

	command := &cobra.Command{
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
		Example:               example,
		Long:                  mapper.Description(cmd),
		RunE:                  handlerFromCommand(ctx, cmd),
		Short:                 mapper.ShortDescription(cmd),
		SilenceErrors:         true,
		SilenceUsage:          true,
		Use:                   commandName + mapper.Usage(cmd),
	}

	hooks := mapper.Hooks(cmd)
	setCobraHooksFromCLIHooks(ctx, command, hooks)

	if err := hooks.BeforeFlagsDefinition(ctx); err != nil {
		return nil, fmt.Errorf("pre-flag-definition hook failed: %w", err)
	}

	if err := setCobraFlagsFromCLIFlags(command.Flags(), mapper.Flags(cmd), options); err != nil {
		return nil, fmt.Errorf("cobra flags build failed: %w", err)
	}

	if err := setCobraFlagsFromCLIFlags(command.PersistentFlags(), mapper.PersistentFlags(cmd), options); err != nil {
		return nil, fmt.Errorf("cobra flags build failed: %w", err)
	}

	return command, nil
}

func getArgs(command *cobra.Command, args []string) ([]string, []string) {
	argsSeparatedAt := command.ArgsLenAtDash()

	switch {
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

func handlerFromCommand(ctx context.Context, cmd cli.Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		args, dashedArgs := getArgs(c, args)

		err := cmd.Execute(ctx, args, dashedArgs)

		var showHelpErr cli.ShowHelpError
		if errors.As(err, &showHelpErr) {
			if showHelpErr.ShowHelp() {
				_ = c.Usage()
			}
			err = errors.Unwrap(err)
		}

		return err
	}
}

func setCobraHooksFromCLIHooks(ctx context.Context, c *cobra.Command, hooks *cli.Hooks) {
	c.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PersistentPreRunE != nil {
			if err := parent.PersistentPreRunE(parent, args); err != nil {
				return err
			}
		}
		return hooks.PersistentBeforeCommandExecution(ctx)
	}
	c.PersistentPostRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PersistentPostRunE != nil {
			if err := parent.PersistentPostRunE(parent, args); err != nil {
				return err
			}
		}
		return hooks.PersistentAfterCommandExecution(ctx)
	}
	c.PreRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PreRunE != nil {
			if err := parent.PreRunE(parent, args); err != nil {
				return err
			}
		}
		return hooks.BeforeCommandExecution(ctx)
	}
	c.PostRunE = func(c *cobra.Command, args []string) error {
		if parent := c.Parent(); parent != nil && parent.PostRunE != nil {
			if err := parent.PostRunE(parent, args); err != nil {
				return err
			}
		}
		return hooks.AfterCommandExecution(ctx)
	}
}

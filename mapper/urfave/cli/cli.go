package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	urfavecli "github.com/urfave/cli/v2"
	"go.uber.org/multierr"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
	"github.com/krostar/cli/mapper"
)

func Execute(ctx context.Context, c *cli.CLI, args []string, opts ...Option) error {
	options := newOptions()
	for _, o := range opts {
		o(options)
	}

	command, err := buildCommandRecursively(ctx, c, options)
	if err != nil {
		return err
	}

	app := urfavecli.App{
		Action:                 command.Action,
		After:                  command.After,
		ArgsUsage:              command.ArgsUsage,
		BashComplete:           command.BashComplete,
		Before:                 command.Before,
		Commands:               command.Subcommands,
		Compiled:               app.BuiltAt(),
		Description:            command.Description,
		ExitErrHandler:         onError,
		Flags:                  command.Flags,
		HelpName:               command.HelpName,
		HideHelp:               command.HideHelp,
		HideHelpCommand:        command.HideHelpCommand,
		Name:                   app.Name(),
		OnUsageError:           command.OnUsageError,
		Usage:                  command.Usage,
		UsageText:              command.UsageText,
		UseShortOptionHandling: command.UseShortOptionHandling,
		Version:                app.Version(),
	}

	return app.RunContext(ctx, args)
}

func buildCommandRecursively(ctx context.Context, cli *cli.CLI, options *options) (*urfavecli.Command, error) {
	ctx = mapper.Context(cli.Command, ctx)

	command, err := buildCommand(ctx, cli.Name, cli.Command, options)
	if err != nil {
		return nil, fmt.Errorf("unable to build urfave/cli command %s: %w", cli.Name, err)
	}

	for _, sub := range cli.SubCommands {
		subCommand, err := buildCommandRecursively(ctx, sub, options)
		if err != nil {
			return nil, fmt.Errorf("unable to build urfave/cli sub-command %s of command %s: %w", sub.Name, cli.Name, err)
		}
		command.Subcommands = append(command.Subcommands, subCommand)
	}

	return command, nil
}

func buildCommand(ctx context.Context, commandName string, cmd cli.Command, options *options) (*urfavecli.Command, error) {
	description := mapper.Description(cmd)
	if examples := mapper.Examples(cmd); len(examples) > 0 {
		description += "\n\nEXAMPLES:\n   " + strings.Join(mapper.Examples(cmd), "\n   ")
	}

	var usage string
	if u := mapper.Usage(cmd); u != "" {
		usage = commandName + u
	}

	hooks := mapper.Hooks(cmd)
	command := &urfavecli.Command{
		Action:                 handlerFromCommand(cmd, hooks),
		After:                  handlerFromHook(hooks.PersistentAfterCommandExecution),
		Before:                 handlerFromHook(hooks.PersistentBeforeCommandExecution),
		Description:            description,
		Name:                   commandName,
		Usage:                  mapper.ShortDescription(cmd),
		UsageText:              usage,
		UseShortOptionHandling: true,
	}

	if err := hooks.BeforeFlagsDefinition(ctx); err != nil {
		return nil, fmt.Errorf("pre-flag-definition hook failed: %w", err)
	}

	var err error
	if command.Flags, err = buildFlags(cmd, options); err != nil {
		return nil, fmt.Errorf("urfave/cli flags build failed: %w", err)
	}

	return command, nil
}

func handlerFromHook(handler func(context.Context) error) func(*urfavecli.Context) error {
	return func(c *urfavecli.Context) error { return handler(c.Context) }
}

func handlerFromCommand(cmd cli.Command, hooks *cli.Hooks) func(*urfavecli.Context) error {
	return func(c *urfavecli.Context) error {
		ctx := c.Context

		if err := hooks.BeforeCommandExecution(ctx); err != nil {
			return err
		}

		args, dashedArgs := getDifferentsArgs(c.Args().Slice())

		var err error
		if handlerErr := cmd.Execute(ctx, args, dashedArgs); handlerErr != nil {
			var errHelp cli.ShowHelpError
			if errors.As(err, &errHelp) {
				handlerErr = flag.ErrHelp
			}
			err = multierr.Append(err, handlerErr)
		}

		if hookErr := hooks.AfterCommandExecution(ctx); hookErr != nil {
			err = multierr.Append(err, hookErr)
		}

		return err
	}
}

func getDifferentsArgs(args []string) ([]string, []string) {
	var (
		beforeDashArgs []string
		afterDashArgs  []string
		foundDash      bool
	)

	for _, arg := range args {
		if foundDash {
			afterDashArgs = append(afterDashArgs, arg)
			continue
		}
		if arg == "--" {
			foundDash = true
			continue
		}
		beforeDashArgs = append(beforeDashArgs, arg)
	}

	return beforeDashArgs, afterDashArgs
}

func onError(c *urfavecli.Context, err error) {
	var errHelp cli.ShowHelpError
	if errors.As(err, &errHelp) {
		_ = urfavecli.ShowCommandHelp(c, c.Command.Name)
	}
}

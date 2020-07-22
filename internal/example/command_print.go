package example

import (
	"context"
	"errors"

	"github.com/krostar/cli"
	"github.com/krostar/logger"
)

type CommandPrint struct {
	logger logger.Logger
}

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

func (cmd *CommandPrint) Hooks() *cli.Hooks {
	return &cli.Hooks{
		BeforeCommandExecution: func(ctx context.Context) error {
			logger, err := getLogger(ctx)
			if err != nil {
				return err
			}
			cmd.logger = logger
			return nil
		},
		BeforeFlagsDefinition: func(context.Context) error {
			return nil
		},
	}
}

func (cmd CommandPrint) Execute(_ context.Context, args []string, dashedArgs []string) error {
	cmd.logger.Info("print command")

	if len(args) == 0 {
		return cli.ErrorShowHelp(errors.New("there should be at least 1 arg to print"))
	}
	if len(args) > 3 {
		return errors.New("there should be no more than 3 arg to print")
	}

	for _, arg := range args {
		cmd.logger.WithField("type", "argument").Infof("%q", arg)
	}
	for _, arg := range dashedArgs {
		cmd.logger.WithField("type", "dashed argument").Infof("%q", arg)
	}
	return nil
}

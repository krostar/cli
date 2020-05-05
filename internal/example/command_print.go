package example

import (
	"context"

	"github.com/krostar/cli"
	"github.com/krostar/logger"
)

type CommandPrint struct {
	logger logger.Logger
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

	for _, arg := range args {
		cmd.logger.WithField("type", "argument").Infof("%q", arg)
	}
	for _, arg := range dashedArgs {
		cmd.logger.WithField("type", "dashed argument").Infof("%q", arg)
	}
	return nil
}

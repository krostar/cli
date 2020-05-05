package example

import (
	"context"
	"fmt"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
	"github.com/krostar/config"
	"github.com/krostar/logger"
	"github.com/krostar/logger/zap"
)

type CommandRoot struct {
	logger logger.Logger
	cfg    commandRootConfig
}

type commandRootConfig struct {
	Logger logger.Config
}

func (CommandRoot) Description() string {
	return app.Name() + ` is a cli app built with abstraction of any cli backend.
The main goal of this app is to demonstrate how easy it is to switch from one cli backend to another.`
}

func (cmd *CommandRoot) PersistentFlags() []cli.Flag {
	return []cli.Flag{
		cli.FlagString("log-format", "", &cmd.cfg.Logger.Formatter, "format the log will be written to"),
		cli.FlagString("log-verbosity", "", &cmd.cfg.Logger.Verbosity, "verbosity level of logs"),
		cli.FlagString("log-output", "", &cmd.cfg.Logger.Output, "output where logs will be written to"),
	}
}

func (cmd *CommandRoot) Hooks() *cli.Hooks {
	var flushLogs func() error

	return &cli.Hooks{
		BeforeFlagsDefinition: func(ctx context.Context) error {
			if err := config.SetDefault(&cmd.cfg); err != nil {
				return fmt.Errorf("unable to set config defaults")
			}
			return nil
		},
		PersistentBeforeCommandExecution: func(ctx context.Context) error {
			logger, flush, err := zap.New(zap.WithConfig(cmd.cfg.Logger))
			if err != nil {
				return fmt.Errorf("unable to build logger: %v", err)
			}

			flushLogs = flush
			setLogger(ctx, logger)

			return nil
		},
		PersistentAfterCommandExecution: func(ctx context.Context) error {
			_ = flushLogs()
			return nil
		},
	}
}

func (cmd CommandRoot) Execute(_ context.Context, args []string, dashedArgs []string) error {
	cmd.logger.Info("root command")
	return cli.ErrorWithExitStatus(cli.ErrShowHelp, 0)
}

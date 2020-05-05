package example

import (
	"context"
	"errors"

	"github.com/krostar/cli"
	"github.com/krostar/logger"
)

type ctxKey uint8

const ctxKeyLogger ctxKey = iota

func setLogger(ctx context.Context, logger logger.Logger) {
	cli.SetMetadata(ctx, ctxKeyLogger, logger)
}

func getLogger(ctx context.Context) (logger.Logger, error) {
	if logger, ok := cli.GetMetadata(ctx, ctxKeyLogger).(logger.Logger); ok {
		return logger, nil
	}
	return nil, errors.New("logger not injected in metadata")
}

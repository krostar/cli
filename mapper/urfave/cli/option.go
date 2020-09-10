package cli

import (
	"fmt"

	urfavecli "github.com/urfave/cli/v2"

	"github.com/krostar/cli"
)

type Option func(o *options)

func WithCustomFlagMapper(f CustomFlagMapperFunc) Option {
	return func(o *options) {
		o.customFlagMapper = f
	}
}

type CustomFlagMapperFunc func(flag cli.Flag) (urfavecli.Flag, error)

type options struct {
	customFlagMapper CustomFlagMapperFunc
}

func newOptions() *options {
	return &options{
		customFlagMapper: func(flag cli.Flag) (urfavecli.Flag, error) {
			return nil, fmt.Errorf("unhandled flag type: %T", flag.Destination())
		},
	}
}

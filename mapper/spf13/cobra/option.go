package cobra

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/krostar/cli"
)

type Option func(o *options)

func WithCustomFlagMapper(f CustomFlagMapperFunc) Option {
	return func(o *options) {
		o.customFlagMapper = f
	}
}

type CustomFlagMapperFunc func(flag cli.Flag, set *pflag.FlagSet) error

type options struct {
	customFlagMapper CustomFlagMapperFunc
}

func newOptions() *options {
	return &options{
		customFlagMapper: func(flag cli.Flag, set *pflag.FlagSet) error {
			return fmt.Errorf("unhandled flag type: %T", flag.Destination())
		},
	}
}

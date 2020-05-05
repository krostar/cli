package cli

import (
	"flag"
	"fmt"
	"time"

	urfavecli "github.com/urfave/cli/v2"
	"go.uber.org/multierr"

	"github.com/krostar/cli"
	"github.com/krostar/cli/mapper"
)

func buildFlags(cmd cli.Command) ([]urfavecli.Flag, error) {
	var (
		flags []urfavecli.Flag
		err   error
	)

	addFlag := func(flag cli.Flag) {
		var alias []string
		if shorthand := flag.Shorthand(); shorthand != "" {
			alias = []string{shorthand}
		}

		switch dest := flag.Destination().(type) {
		case *bool:
			flags = append(flags, &urfavecli.BoolFlag{Name: flag.Name(), Aliases: alias, Usage: flag.Description(), Value: *dest, Destination: dest})
		case *[]bool:
			flags = append(flags, newFlagBoolSlice(flag, dest))
		case *string:
			flags = append(flags, &urfavecli.StringFlag{Name: flag.Name(), Aliases: alias, Usage: flag.Description(), Value: *dest, Destination: dest})
		case *[]string:
			flags = append(flags, newFlagStringSlice(flag, dest))
		case *int:
			flags = append(flags, &urfavecli.IntFlag{Name: flag.Name(), Aliases: alias, Usage: flag.Description(), Value: *dest, Destination: dest})
		case *[]int:
			flags = append(flags, newFlagIntSlice(flag, dest))
		case *uint:
			flags = append(flags, &urfavecli.UintFlag{Name: flag.Name(), Aliases: alias, Usage: flag.Description(), Value: *dest, Destination: dest})
		case *[]uint:
			flags = append(flags, newFlagUintSlice(flag, dest))
		case *float32:
			flags = append(flags, newFlagFloat32(flag, dest))
		case *[]float32:
			flags = append(flags, newFlagFloat32Slice(flag, dest))
		case *float64:
			flags = append(flags, &urfavecli.Float64Flag{Name: flag.Name(), Aliases: alias, Usage: flag.Description(), Value: *dest, Destination: dest})
		case *[]float64:
			flags = append(flags, newFlagFloat64Slice(flag, dest))
		case *time.Duration:
			flags = append(flags, &urfavecli.DurationFlag{Name: flag.Name(), Aliases: alias, Usage: flag.Description(), Value: *dest, Destination: dest})
		case *[]time.Duration:
			flags = append(flags, newFlagDurationSlice(flag, dest))
		default:
			err = multierr.Append(err, fmt.Errorf("unhandled flag type: %T", dest))
		}
	}

	for _, flag := range mapper.PersistentFlags(cmd) {
		addFlag(flag)
	}

	for _, flag := range mapper.Flags(cmd) {
		addFlag(flag)
	}

	return flags, err
}

type customFlag struct {
	flag        cli.Flag
	destination flag.Value

	Usage       string
	DefaultText string
	Value       interface{}
}

func newCustomFlag(flag cli.Flag, value interface{}, destination flag.Value) *customFlag {
	return &customFlag{
		flag:        flag,
		destination: destination,
		Usage:       flag.Description(),
		DefaultText: destination.String(),
		Value:       value,
	}
}

func (f customFlag) Names() []string {
	names := []string{f.flag.Name()}
	if shorthand := f.flag.Shorthand(); shorthand != "" {
		names = append(names, shorthand)
	}
	return names
}

func (f customFlag) Apply(set *flag.FlagSet) error {
	for _, name := range f.Names() {
		set.Var(f.destination, name, f.flag.Description())
	}
	return nil
}

type customValue struct {
	isSet bool
}

func (v customValue) Serialize() string {
	return "e3Y-RVt=h9!A9/Wg}TmhEwp@?|*)k_g;4oxLA'ir]AY*&wKndoCbY3yxhUNgePii.TrtJ%*y#u.zTcH:'(k*<:Vdqn#4-yuW@=NN"
}

func (v customValue) IsSet() bool { return v.isSet }

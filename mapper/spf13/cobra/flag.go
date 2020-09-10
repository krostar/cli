package cobra

import (
	"time"

	"github.com/spf13/pflag"

	"github.com/krostar/cli"
)

func buildFlags(set *pflag.FlagSet, flags []cli.Flag, options *options) error {
	var err error

	for _, flag := range flags {
		switch dest := flag.Destination().(type) {
		case *bool:
			set.BoolVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]bool:
			set.BoolSliceVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *string:
			set.StringVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]string:
			set.StringSliceVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *int:
			set.IntVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]int:
			set.IntSliceVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *uint:
			set.UintVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]uint:
			set.UintSliceVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *float32:
			set.Float32VarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]float32:
			set.VarP(newFlagFloat32Value(*dest, dest), flag.Name(), flag.Shorthand(), flag.Description())
		case *float64:
			set.Float64VarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]float64:
			set.VarP(newFlagFloat64Value(*dest, dest), flag.Name(), flag.Shorthand(), flag.Description())
		case *time.Duration:
			set.DurationVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		case *[]time.Duration:
			set.DurationSliceVarP(dest, flag.Name(), flag.Shorthand(), *dest, flag.Description())
		default:
			err = options.customFlagMapper(flag, set)
		}
	}

	return err
}

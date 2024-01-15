package spf13cobra

import (
	"github.com/spf13/pflag"

	"github.com/krostar/cli"
)

func setCobraFlagsFromCLIFlags(set *pflag.FlagSet, flags []cli.Flag) {
	for _, flag := range flags {
		set.VarP(&flagValuer{flag}, flag.LongName(), flag.ShortName(), flag.Description())
	}
}

type flagValuer struct{ cli.FlagValuer }

func (flag *flagValuer) Set(raw string) error {
	return flag.FlagValuer.FromString(raw)
}

func (flag *flagValuer) Type() string {
	return flag.FlagValuer.TypeRepr()
}

func (flag *flagValuer) String() string { return flag.FlagValuer.String() }

package cli

import (
	"strings"

	"github.com/krostar/cli"
)

type flagStringSlice struct {
	*CustomFlag
	*flagStringSliceValue
}

func newFlagStringSlice(flag cli.Flag, dest *[]string) *flagStringSlice {
	v := &flagStringSliceValue{dest: dest}
	return &flagStringSlice{
		CustomFlag:           NewCustomFlag(flag, *dest, v),
		flagStringSliceValue: v,
	}
}

type flagStringSliceValue struct {
	customValue
	dest *[]string
}

func (v flagStringSliceValue) String() string { return "[" + strings.Join(*v.dest, ",") + "]" }
func (v *flagStringSliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []string{}
	}

	*v.dest = append(*v.dest, value)
	v.isSet = true

	return nil
}

package cli

import (
	"fmt"
	"strconv"
	"strings"

	urfavecli "github.com/urfave/cli/v2"

	"github.com/krostar/cli"
)

func newFlagBoolSlice(flag cli.Flag, dest *[]bool) *flagBoolSlice {
	v := &flagBoolSliceValue{dest: dest}
	return &flagBoolSlice{
		customFlag:         newCustomFlag(flag, *dest, v),
		flagBoolSliceValue: v,
	}
}

type flagBoolSlice struct {
	*customFlag
	*flagBoolSliceValue
}

func (f flagBoolSlice) String() string { return urfavecli.FlagStringer(f) }

type flagBoolSliceValue struct {
	customValue
	dest *[]bool
}

func (v flagBoolSliceValue) String() string {
	if v.dest == nil || *v.dest == nil {
		return ""
	}

	var repr []string
	for _, d := range *v.dest {
		repr = append(repr, strconv.FormatBool(d))
	}
	return "[" + strings.Join(repr, ",") + "]"
}

func (v *flagBoolSliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []bool{}
	}

	vb, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("unable to parse bool: %v", err)
	}

	*v.dest = append(*v.dest, vb)
	v.isSet = true

	return nil
}

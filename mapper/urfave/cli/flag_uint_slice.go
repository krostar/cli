package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/krostar/cli"
)

type flagUintSlice struct {
	*CustomFlag
	*flagUintSliceValue
}

func newFlagUintSlice(flag cli.Flag, dest *[]uint) *flagUintSlice {
	v := &flagUintSliceValue{dest: dest}
	return &flagUintSlice{
		CustomFlag:         NewCustomFlag(flag, *dest, v),
		flagUintSliceValue: v,
	}
}

type flagUintSliceValue struct {
	customValue
	dest *[]uint
}

func (v flagUintSliceValue) String() string {
	var repr []string
	for _, nb := range *v.dest {
		repr = append(repr, strconv.FormatUint(uint64(nb), 64))
	}
	return "[" + strings.Join(repr, ",") + "]"
}

func (v *flagUintSliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []uint{}
	}

	nb, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse number: %v", err)
	}

	*v.dest = append(*v.dest, uint(nb))
	v.isSet = true

	return nil
}

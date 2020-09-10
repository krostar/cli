package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/krostar/cli"
)

type flagDurationSlice struct {
	*CustomFlag
	*flagDurationSliceValue
}

func newFlagDurationSlice(flag cli.Flag, dest *[]time.Duration) *flagDurationSlice {
	v := &flagDurationSliceValue{dest: dest}
	return &flagDurationSlice{
		CustomFlag:             NewCustomFlag(flag, *dest, v),
		flagDurationSliceValue: v,
	}
}

type flagDurationSliceValue struct {
	customValue
	dest *[]time.Duration
}

func (v flagDurationSliceValue) String() string {
	var repr []string
	for _, d := range *v.dest {
		repr = append(repr, d.String())
	}
	return "[" + strings.Join(repr, ",") + "]"
}

func (v *flagDurationSliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []time.Duration{}
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("unable to parse number: %v", err)
	}

	*v.dest = append(*v.dest, d)
	v.isSet = true

	return nil
}

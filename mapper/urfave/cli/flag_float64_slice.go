package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/krostar/cli"
)

type flagFloat64Slice struct {
	*CustomFlag
	*flagFloat64SliceValue
}

func newFlagFloat64Slice(flag cli.Flag, dest *[]float64) *flagFloat64Slice {
	v := &flagFloat64SliceValue{dest: dest}
	return &flagFloat64Slice{
		CustomFlag:            NewCustomFlag(flag, *dest, v),
		flagFloat64SliceValue: v,
	}
}

type flagFloat64SliceValue struct {
	customValue
	dest *[]float64
}

func (v flagFloat64SliceValue) String() string {
	var repr []string
	for _, nb := range *v.dest {
		repr = append(repr, strconv.FormatFloat(nb, 'f', -1, 64))
	}
	return "[" + strings.Join(repr, ",") + "]"
}

func (v *flagFloat64SliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []float64{}
	}

	nb, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("unable to parse number: %v", err)
	}

	*v.dest = append(*v.dest, nb)
	v.isSet = true

	return nil
}

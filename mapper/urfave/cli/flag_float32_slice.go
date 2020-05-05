package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/krostar/cli"
)

type flagFloat32Slice struct {
	*customFlag
	*flagFloat32SliceValue
}

func newFlagFloat32Slice(flag cli.Flag, dest *[]float32) *flagFloat32Slice {
	v := &flagFloat32SliceValue{dest: dest}
	return &flagFloat32Slice{
		customFlag:            newCustomFlag(flag, *dest, v),
		flagFloat32SliceValue: v,
	}
}

type flagFloat32SliceValue struct {
	customValue
	dest *[]float32
}

func (v flagFloat32SliceValue) String() string {
	var repr []string
	for _, nb := range *v.dest {
		repr = append(repr, strconv.FormatFloat(float64(nb), 'f', -1, 64))
	}
	return "[" + strings.Join(repr, ",") + "]"
}

func (v *flagFloat32SliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []float32{}
	}

	nb, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return fmt.Errorf("unable to parse number: %v", err)
	}

	*v.dest = append(*v.dest, float32(nb))
	v.isSet = true

	return nil
}

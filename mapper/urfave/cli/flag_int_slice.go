package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/krostar/cli"
)

type flagIntSlice struct {
	*customFlag
	*flagIntSliceValue
}

func newFlagIntSlice(flag cli.Flag, dest *[]int) *flagIntSlice {
	v := &flagIntSliceValue{dest: dest}
	return &flagIntSlice{
		customFlag:        newCustomFlag(flag, *dest, v),
		flagIntSliceValue: v,
	}
}

type flagIntSliceValue struct {
	customValue
	dest *[]int
}

func (v flagIntSliceValue) String() string {
	var repr []string
	for _, nb := range *v.dest {
		repr = append(repr, strconv.Itoa(nb))
	}
	return "[" + strings.Join(repr, ",") + "]"
}

func (v *flagIntSliceValue) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	if !v.isSet {
		*v.dest = []int{}
	}

	nb, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("unable to parse number: %v", err)
	}

	*v.dest = append(*v.dest, nb)
	v.isSet = true

	return nil
}

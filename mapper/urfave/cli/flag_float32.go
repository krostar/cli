package cli

import (
	"fmt"
	"strconv"

	urfavecli "github.com/urfave/cli/v2"

	"github.com/krostar/cli"
)

func newFlagFloat32(flag cli.Flag, dest *float32) *flagFloat32 {
	v := &flagFloat32Value{dest: dest}
	return &flagFloat32{
		CustomFlag:       NewCustomFlag(flag, *dest, v),
		flagFloat32Value: v,
	}
}

type flagFloat32 struct {
	*CustomFlag
	*flagFloat32Value
}

func (f flagFloat32) String() string { return urfavecli.FlagStringer(f) }

type flagFloat32Value struct {
	customValue
	dest *float32
}

func (v flagFloat32Value) String() string {
	if v.dest == nil {
		return ""
	}

	return strconv.FormatFloat(float64(*v.dest), 'f', -1, 32)
}

func (v *flagFloat32Value) Set(value string) error {
	if value == v.Serialize() {
		return nil
	}

	vb, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return fmt.Errorf("unable to parse bool: %v", err)
	}

	*v.dest = float32(vb)
	v.isSet = true

	return nil
}

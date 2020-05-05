package cobra

import (
	"fmt"
	"strconv"
	"strings"
)

type flagFloat32SliceValue struct {
	value   *[]float32
	changed bool
}

func newFlagFloat32Value(val []float32, p *[]float32) *flagFloat32SliceValue {
	fsv := new(flagFloat32SliceValue)
	fsv.value = p
	*fsv.value = val
	return fsv
}

func (s *flagFloat32SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]float32, len(ss))
	for i, rawF := range ss {
		f, err := strconv.ParseFloat(rawF, 32)
		if err != nil {
			return err
		}
		out[i] = float32(f)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}

func (s *flagFloat32SliceValue) Type() string {
	return "float32Slice"
}

func (s *flagFloat32SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, f := range *s.value {
		out[i] = fmt.Sprintf("%f", f)
	}
	return "[" + strings.Join(out, ",") + "]"
}

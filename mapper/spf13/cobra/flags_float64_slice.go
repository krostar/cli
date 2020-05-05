package cobra

import (
	"fmt"
	"strconv"
	"strings"
)

type flagFloat64SliceValue struct {
	value   *[]float64
	changed bool
}

func newFlagFloat64Value(val []float64, p *[]float64) *flagFloat64SliceValue {
	fsv := new(flagFloat64SliceValue)
	fsv.value = p
	*fsv.value = val
	return fsv
}

func (s *flagFloat64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, rawF := range ss {
		var err error
		out[i], err = strconv.ParseFloat(rawF, 64)
		if err != nil {
			return err
		}
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}

func (s *flagFloat64SliceValue) Type() string {
	return "float64Slice"
}

func (s *flagFloat64SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, f := range *s.value {
		out[i] = fmt.Sprintf("%f", f)
	}
	return "[" + strings.Join(out, ",") + "]"
}

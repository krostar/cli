package cli

import (
	"fmt"
	"strings"
)

// FlagValuer defines how to set and get the value of a flagValue.
type FlagValuer interface {
	// TypeRepr returns a representation of the underlying type of the flagValue. Example: 'int'.
	TypeRepr() string
	// ToString returns a text representation of the flagValue's value.
	ToString() (string, error)
	// FromString parses and set the provided flagValue value.
	FromString(str string) error
}

// NewCustomFlagValuer creates a flag valuer from a parse function returning somethind that implements stringer.
func NewCustomFlagValuer[T fmt.Stringer](destination *T, parse func(string) (T, error)) FlagValuer {
	return &flagValuer[T]{
		value: destination,
		parse: parse,
	}
}

type flagValuer[T fmt.Stringer] struct {
	value *T
	parse func(string) (T, error)
}

func (v *flagValuer[T]) TypeRepr() string {
	return fmt.Sprintf("%T", *v.value)
}

func (v *flagValuer[T]) ToString() (string, error) {
	return (*v.value).String(), nil
}

func (v *flagValuer[T]) FromString(raw string) error {
	value, err := v.parse(strings.TrimSpace(raw))
	if err != nil {
		return err
	}

	*v.value = value

	return nil
}

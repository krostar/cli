package cli

import (
	"fmt"
	"strings"
)

// FlagValuer defines how to set and get the value of a flag.
type FlagValuer interface {
	// FromString parses and set the value.
	FromString(str string) error
	// IsSet returns whether the flag value has been set.
	IsSet() bool
	// String returns a representation of the value.
	String() string
	// TypeRepr returns a representation of the underlying type of the value. Example: 'int'.
	TypeRepr() string
}

// NewFlagValuer creates a flag valuer from a parse and toString functions. It panics if any parameters are nil.
func NewFlagValuer[T any](destination *T, parse func(string) (T, error), toString func(T) string) FlagValuer {
	if destination == nil {
		panic("destination is nil")
	}

	if parse == nil {
		panic("parse is nil")
	}

	if toString == nil {
		panic("toString is nil")
	}

	return &flagValuer[T]{
		value:    destination,
		parse:    parse,
		toString: toString,
	}
}

// NewStringerFlagValuer creates a flag valuer from a parse function returning something that implements stringer.
// See NewFlagValuer for more details.
func NewStringerFlagValuer[T fmt.Stringer](destination *T, parse func(string) (T, error)) FlagValuer {
	return NewFlagValuer(destination, parse, func(t T) string { return t.String() })
}

type flagValuer[T any] struct {
	value    *T
	changed  bool
	parse    func(string) (T, error)
	toString func(T) string
}

func (v *flagValuer[T]) FromString(raw string) error {
	value, err := v.parse(strings.TrimSpace(raw))
	if err != nil {
		return err
	}

	*v.value = value
	v.changed = true

	return nil
}

func (v *flagValuer[T]) IsSet() bool { return v.changed }

func (v *flagValuer[T]) String() string { return v.toString(*v.value) }

func (v *flagValuer[T]) TypeRepr() string { return fmt.Sprintf("%T", *v.value) }

package cli

import (
	"fmt"
	"strings"
)

// FlagValuer defines the interface for getting and setting the value of a flag.
// It abstracts the underlying type of the flag, allowing for type-safe flag handling
// while providing a common interface for CLI frameworks to interact with.
//
// By implementing this interface, you can create custom flag types that work
// seamlessly with the CLI framework.
type FlagValuer interface {
	// Destination returns a pointer to the variable that stores the flag's value.
	// This allows CLI frameworks to access the underlying value directly.
	Destination() any

	// FromString parses a string and sets the flag's value.
	// This is called when a flag is specified on the command line.
	// It should convert the string representation to the appropriate type
	// and store it in the destination.
	FromString(str string) error

	// IsSet returns true if the flag has been set (i.e., its value has been
	// modified from the default). This allows commands to determine if a
	// flag was explicitly provided by the user.
	IsSet() bool

	// String returns a string representation of the flag's current value.
	// This is used for displaying the flag's value in help text and error messages.
	String() string

	// TypeRepr returns a string representing the underlying type of the
	// flag's value (e.g., "int", "string", "bool"). This is used in help text
	// to indicate the expected type of the flag's value.
	TypeRepr() string
}

// NewFlagValuer creates a new FlagValuer instance for any type.
// It provides a generic implementation of the FlagValuer interface that can
// work with any Go type.
//
// Parameters:
//   - destination: Pointer to where the flag value will be stored
//   - parse: Function to parse a string into the destination type
//   - toString: Function to convert the destination type to a string
//
// Returns a FlagValuer that manages the value at destination.
//
// It panics if any of the parameters is nil.
//
// Example:
//
//	var count int
//	valuer := NewFlagValuer(
//	    &count,
//	    func(s string) (int, error) { return strconv.Atoi(s) },
//	    func(i int) string { return strconv.Itoa(i) },
//	)
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

// NewStringerFlagValuer creates a FlagValuer for types that implement the fmt.Stringer interface.
// It simplifies the creation of FlagValuers for types that already have a String() method.
//
// Parameters:
//   - destination: Pointer to where the flag value will be stored
//   - parse: Function to parse a string into the destination type
//
// Returns a FlagValuer that manages the value at destination, using the type's
// String() method for string representation.
//
// Example:
//
//	var duration time.Duration
//	valuer := NewStringerFlagValuer(
//	    &duration,
//	    func(s string) (time.Duration, error) { return time.ParseDuration(s) },
//	)
func NewStringerFlagValuer[T fmt.Stringer](destination *T, parse func(string) (T, error)) FlagValuer {
	return NewFlagValuer(destination, parse, func(t T) string { return t.String() })
}

// flagValuer is a generic implementation of the FlagValuer interface.
// It handles the conversion between string and typed values, and tracks
// whether the flag has been set.
type flagValuer[T any] struct {
	value    *T
	changed  bool
	parse    func(string) (T, error)
	toString func(T) string
}

// Destination returns the pointer to the flag's value.
func (v *flagValuer[T]) Destination() any { return v.value }

// FromString parses a string and sets the flag's value.
// It trims whitespace from the input string, parses it using the provided
// parse function, and updates the destination value.
func (v *flagValuer[T]) FromString(raw string) error {
	value, err := v.parse(strings.TrimSpace(raw))
	if err != nil {
		return err
	}

	*v.value = value
	v.changed = true

	return nil
}

// IsSet returns true if the flag value has been set.
func (v *flagValuer[T]) IsSet() bool { return v.changed }

// String returns the string representation of the flag's value.
func (v *flagValuer[T]) String() string { return v.toString(*v.value) }

// TypeRepr returns a string representation of the flag's type.
func (v *flagValuer[T]) TypeRepr() string { return fmt.Sprintf("%T", *v.value) }

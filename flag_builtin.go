package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// NewBuiltinFlag creates a Flag for built-in types (int, string, bool, etc.).
// It handles the conversion between string representations and the underlying
// Go type. See NewFlag for more details.
func NewBuiltinFlag[T builtins](longName, shortName string, destination *T, description string) Flag {
	return NewFlag(
		longName, shortName,
		NewFlagValuer(destination, builtinFromString[T], builtinToString[T]),
		description,
	)
}

// NewBuiltinPointerFlag creates a Flag for pointers to built-in types.
// See NewBuiltinFlag for more details.
func NewBuiltinPointerFlag[T builtins](longName, shortName string, destination **T, description string) Flag {
	return NewFlag(
		longName, shortName,
		NewFlagValuer(destination,
			func(s string) (*T, error) {
				b, err := builtinFromString[T](s)
				if err != nil {
					return nil, err
				}
				return &b, nil
			},
			func(t *T) string {
				if t != nil {
					return builtinToString[T](*t)
				}
				return "<nil>"
			},
		), description,
	)
}

// NewBuiltinSliceFlag creates a Flag for slices of built-in types.
// The flag value is expected to be a comma-separated list of values.
// See NewBuiltinFlag for more details.
func NewBuiltinSliceFlag[T builtins](longName, shortName string, destination *[]T, description string) Flag {
	return NewFlag(
		longName, shortName,
		NewFlagValuer(destination,
			func(raw string) ([]T, error) {
				rawValues := strings.Split(raw, ",")
				values := make([]T, len(rawValues))

				for i, rawValue := range rawValues {
					value, err := builtinFromString[T](strings.TrimSpace(rawValue))
					if err != nil {
						return nil, err
					}

					values[i] = value
				}

				return values, nil
			},
			func(values []T) string {
				valuesRepr := make([]string, len(values))
				for i, value := range values {
					valuesRepr[i] = builtinToString(value)
				}
				return "[" + strings.Join(valuesRepr, ",") + "]"
			},
		), description,
	)
}

type builtins interface {
	bool | string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | complex64 | complex128
}

// builtinFromString converts a string to a built-in type. It uses a generic
// type parameter constrained by `builtins` and a type switch to handle the
// different conversions.
//
//nolint:revive,errcheck // unchecked-type-assertion: linter does not like the (T) cast that is unchecked, but nothing to worry about here
func builtinFromString[T builtins](raw string) (T, error) {
	newT := *new(T)

	switch t := any(newT).(type) {
	case bool:
		v, err := strconv.ParseBool(raw)
		return any(v).(T), err
	case string:
		v := raw
		return any(v).(T), nil
	case int:
		v, err := strconv.ParseInt(raw, 10, 0)
		return any(int(v)).(T), err
	case int8:
		v, err := strconv.ParseInt(raw, 10, 8)
		return any(int8(v)).(T), err
	case int16:
		v, err := strconv.ParseInt(raw, 10, 16)
		return any(int16(v)).(T), err
	case int32:
		v, err := strconv.ParseInt(raw, 10, 32)
		return any(int32(v)).(T), err
	case int64:
		v, err := strconv.ParseInt(raw, 10, 64)
		return any(v).(T), err
	case uint:
		v, err := strconv.ParseUint(raw, 10, 0)
		return any(uint(v)).(T), err
	case uint8:
		v, err := strconv.ParseUint(raw, 10, 8)
		return any(uint8(v)).(T), err
	case uint16:
		v, err := strconv.ParseUint(raw, 10, 16)
		return any(uint16(v)).(T), err
	case uint32:
		v, err := strconv.ParseUint(raw, 10, 32)
		return any(uint32(v)).(T), err
	case uint64:
		v, err := strconv.ParseUint(raw, 10, 64)
		return any(v).(T), err
	case float32:
		v, err := strconv.ParseFloat(raw, 32)
		return any(float32(v)).(T), err
	case float64:
		v, err := strconv.ParseFloat(raw, 64)
		return any(v).(T), err
	case complex64:
		v, err := strconv.ParseComplex(raw, 64)
		return any(complex64(v)).(T), err
	case complex128:
		v, err := strconv.ParseComplex(raw, 128)
		return any(v).(T), err
	default:
		return newT, fmt.Errorf("unhandled type %T", t)
	}
}

// builtinToString converts a built-in type to its string representation.
//
//nolint:mnd // don't lint for hardcoded number for precision
func builtinToString[T builtins](t T) string {
	switch t := any(t).(type) {
	case bool:
		return strconv.FormatBool(t)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case int8:
		return strconv.FormatInt(int64(t), 10)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', 4, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', 4, 64)
	case complex64:
		return strconv.FormatComplex(complex128(t), 'f', 4, 64)
	case complex128:
		return strconv.FormatComplex(t, 'f', 4, 64)
	case string:
		return t
	case int:
		return strconv.FormatInt(int64(t), 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	default:
		return ""
	}
}

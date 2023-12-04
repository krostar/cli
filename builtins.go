package cli

import (
	"fmt"
	"strconv"
)

type builtins interface {
	bool | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128 | string | int | uint
}

func builtinFromString[T builtins](raw string) (T, error) {
	newT := *new(T)

	switch t := any(newT).(type) {
	case bool:
		v, err := strconv.ParseBool(raw)
		return any(v).(T), err
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
	case string:
		v := raw
		return any(v).(T), nil
	case int:
		v, err := strconv.ParseInt(raw, 10, 64)
		return any(int(v)).(T), err
	case uint:
		v, err := strconv.ParseUint(raw, 10, 64)
		return any(uint(v)).(T), err
	default:
		return newT, fmt.Errorf("unhandled type %T", t)
	}
}

//nolint:gomnd // don't lint for hardcoded number for precision
func builtinToString[T builtins](t T) (string, error) {
	switch t := any(t).(type) {
	case bool:
		return strconv.FormatBool(t), nil
	case uint8:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint64:
		return strconv.FormatUint(t, 10), nil
	case int8:
		return strconv.FormatInt(int64(t), 10), nil
	case int16:
		return strconv.FormatInt(int64(t), 10), nil
	case int32:
		return strconv.FormatInt(int64(t), 10), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', 4, 32), nil
	case float64:
		return strconv.FormatFloat(t, 'f', 4, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(t), 'f', 4, 64), nil
	case complex128:
		return strconv.FormatComplex(t, 'f', 4, 64), nil
	case string:
		return t, nil
	case int:
		return strconv.FormatInt(int64(t), 10), nil
	case uint:
		return strconv.FormatUint(uint64(t), 10), nil
	default:
		return "", fmt.Errorf("unhandled type %T", t)
	}
}

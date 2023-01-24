package cli

import (
	"fmt"
	"strings"
)

// NewSliceFlag creates a new flagValue which underlying type is any slice of builtin type declared in builtins.go.
//
//	longName is the long flagValue name, like --longname ; cannot be empty.
//	shortName is the short flagValue name ; usually 1 character, like -s ; can be empty.
//	destination is a pointer on the variable on which flagValue's value will be stored ; cannot be nil.
//	description is a short text explaining the flagValue ; can be empty.
func NewSliceFlag[T builtins](longName, shortName string, destination *[]T, description string) Flag {
	return NewCustomFlag(longName, shortName, &flagBuiltinSliceValue[T]{values: destination}, description)
}

type flagBuiltinSliceValue[T builtins] struct {
	values *[]T
}

func (v *flagBuiltinSliceValue[T]) TypeRepr() string {
	return fmt.Sprintf("%T", *v.values)
}

func (v *flagBuiltinSliceValue[T]) ToString() (string, error) {
	valuesRepr := make([]string, len(*v.values))
	for i, value := range *v.values {
		valueRepr, err := builtinToString(value)
		if err != nil {
			return "", err
		}

		valuesRepr[i] = valueRepr
	}

	return "[" + strings.Join(valuesRepr, ",") + "]", nil
}

func (v *flagBuiltinSliceValue[T]) FromString(raw string) error {
	rawValues := strings.Split(raw, ",")
	values := make([]T, len(rawValues))

	for i, rawValue := range rawValues {
		value, err := builtinFromString[T](strings.TrimSpace(rawValue))
		if err != nil {
			return err
		}

		values[i] = value
	}

	*v.values = append(*v.values, values...)

	return nil
}

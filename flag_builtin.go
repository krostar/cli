package cli

import (
	"fmt"
	"strings"
)

// NewFlag creates a new flagValue which underlying type is any builtin type declared in builtins.go.
//
//	longName is the long flagValue name, like --longname ; cannot be empty.
//	shortName is the short flagValue name ; usually 1 character, like -s ; can be empty.
//	destination is a pointer on the variable on which flagValue's value will be stored ; cannot be nil.
//	description is a short text explaining the flagValue ; can be empty.
func NewFlag[T builtins](longName string, shortName string, destination *T, description string) Flag {
	return NewCustomFlag(longName, shortName, &flagBuiltinValue[T]{value: destination}, description)
}

type flagBuiltinValue[T builtins] struct {
	value *T
}

func (v *flagBuiltinValue[T]) TypeRepr() string {
	return fmt.Sprintf("%T", *v.value)
}

func (v *flagBuiltinValue[T]) ToString() (string, error) {
	return builtinToString(*v.value)
}

func (v *flagBuiltinValue[T]) FromString(raw string) error {
	value, err := builtinFromString[T](strings.TrimSpace(raw))
	if err != nil {
		return err
	}

	*v.value = value

	return nil
}

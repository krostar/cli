package clicfgvalidator

import (
	"errors"
	"reflect"
)

// Validate walks cfg and calls Validate() on any type that implements it.
// When a type implements Validate(), its sub-fields are not visited, the type
// is responsible for validating them.
func Validate[T any](cfg *T) error {
	return recursivelyValidate(reflect.ValueOf(cfg))
}

func recursivelyValidate(v reflect.Value) error {
	type validator interface {
		Validate() error
	}

	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return nil
		}

		if vv, ok := v.Interface().(validator); ok {
			return vv.Validate()
		}

		return recursivelyValidate(v.Elem())

	case reflect.Struct:
		if v.CanAddr() {
			if vv, ok := v.Addr().Interface().(validator); ok {
				return vv.Validate()
			}
		}

		var errs []error

		for i := range v.NumField() {
			if v.Type().Field(i).PkgPath != "" {
				continue
			}

			errs = append(errs, recursivelyValidate(v.Field(i)))
		}

		return errors.Join(errs...)

	default:
		return nil
	}
}

package clicfgvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

// Validate walks cfg and calls Validate() on any type that implements it.
// When a type implements Validate(), its sub-fields are not visited, the type
// is responsible for validating them.
// Errors are annotated with the dot-separated field path, e.g. "Sub.Field: invalid".
func Validate[T any](cfg *T) error {
	return recursivelyValidate(reflect.ValueOf(cfg), "")
}

func recursivelyValidate(v reflect.Value, path string) error {
	type validator interface {
		Validate() error
	}

	wrapErr := func(err error) error {
		if err == nil || path == "" {
			return err
		}
		return fmt.Errorf("%s: %w", path, err)
	}

	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return nil
		}

		if vv, ok := v.Interface().(validator); ok {
			return wrapErr(vv.Validate())
		}

		return recursivelyValidate(v.Elem(), path)

	case reflect.Struct:
		if v.CanAddr() {
			if vv, ok := v.Addr().Interface().(validator); ok {
				return wrapErr(vv.Validate())
			}
		}

		var errs []error

		for i := range v.NumField() {
			if v.Type().Field(i).PkgPath != "" {
				continue
			}

			fieldName := v.Type().Field(i).Name
			fieldPath := fieldName

			if path != "" {
				fieldPath = path + "." + fieldName
			}

			errs = append(errs, recursivelyValidate(v.Field(i), fieldPath))
		}

		return errors.Join(errs...)

	default:
		return nil
	}
}

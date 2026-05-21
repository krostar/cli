package sourcedefault

import (
	"context"
	"reflect"

	clicfg "github.com/krostar/cli/cfg"
)

// Source returns a SourceFunc that sets default values for a config.
// It walks the config recursively and calls SetDefault() on any field that implements it.
func Source[T any]() clicfg.SourceFunc[T] {
	return func(_ context.Context, cfg *T) error {
		recursivelySetDefault(reflect.ValueOf(cfg))
		return nil
	}
}

// recursivelySetDefault walks v bottom-up and calls SetDefault() on every node that implements it.
// It returns true if at least one SetDefault was called within the subtree, which drives nil-pointer
// materialization: a nil pointer is only set when something inside actually needs defaults applied.
func recursivelySetDefault(v reflect.Value) bool {
	type defaultSetter interface {
		SetDefault()
	}

	switch v.Kind() {
	case reflect.Pointer:
		tmpV := v

		if v.IsNil() {
			tmpV = reflect.New(v.Type().Elem())
		}

		anySet := recursivelySetDefault(tmpV.Elem())
		if setter, ok := tmpV.Interface().(defaultSetter); ok {
			setter.SetDefault()

			anySet = true
		}

		if v.IsNil() && anySet {
			v.Set(tmpV)
		}

		return anySet

	case reflect.Struct:
		t := v.Type()

		var anySet bool

		for i := range v.NumField() {
			// Unexported fields are never addressable via reflect and cannot implement interfaces.
			if t.Field(i).PkgPath != "" {
				continue
			}

			field := v.Field(i)

			// Recurse first so nested defaults are applied before this field's own SetDefault.
			if recursivelySetDefault(field) {
				anySet = true
			}

			if field.Kind() == reflect.Struct && field.CanAddr() {
				if setter, ok := field.Addr().Interface().(defaultSetter); ok {
					setter.SetDefault()

					anySet = true
				}
			}
		}

		return anySet

	default:
		return false
	}
}

package sourceflag

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/multierr"

	"github.com/krostar/cli"
	clicfg "github.com/krostar/cli/cfg"
)

// Source applies the flag values to the config.
func Source[T any](flagDest *T) clicfg.SourceFunc[T] {
	return func(ctx context.Context, cfg *T) error {
		pointersToValuesSetByFlags := make(map[uintptr]struct{})
		{
			localFlags, persistentFlags := cli.GetInitializedFlagsFromContext(ctx)
			for _, flag := range append(localFlags, persistentFlags...) {
				if flag.IsSet() {
					pointersToValuesSetByFlags[uintptr(reflect.ValueOf(flag.Destination()).UnsafePointer())] = struct{}{}
				}
			}
			if len(pointersToValuesSetByFlags) == 0 {
				return nil
			}
		}

		if err := recursivelyWalkThroughReflectValue(pointersToValuesSetByFlags, reflect.ValueOf(flagDest).Elem(), reflect.ValueOf(cfg).Elem()); err != nil {
			return fmt.Errorf("unable to walk through config: %v", err)
		}

		if len(pointersToValuesSetByFlags) != 0 {
			return errors.New("some values where not find, make sure flag values all points to config")
		}

		return nil
	}
}

func recursivelyWalkThroughReflectValue(pointers map[uintptr]struct{}, v1, v2 reflect.Value) error {
	if len(pointers) == 0 {
		return nil
	}

	switch k := v1.Type().Kind(); k { //nolint:exhaustive // all other cases handled in default
	case reflect.Pointer:
		if v1.IsNil() {
			return nil
		}

		v1ptr := uintptr(v1.Addr().UnsafePointer())
		if _, ok := pointers[v1ptr]; !ok {
			return recursivelyWalkThroughReflectValue(pointers, v1.Elem(), v2.Elem())
		}

		v2.Set(v1)
		delete(pointers, v1ptr)
		return nil

	case reflect.Struct:
		var errs []error
		for i := range v1.NumField() {
			errs = append(errs, recursivelyWalkThroughReflectValue(pointers, v1.Field(i), v2.Field(i)))
		}
		return multierr.Combine(errs...)

	default:
		if !v1.CanAddr() {
			return fmt.Errorf("unable to get address of value: %v", v1)
		}

		v1ptr := uintptr(v1.Addr().UnsafePointer())
		if _, ok := pointers[v1ptr]; !ok {
			return nil
		}

		v2.Set(v1)
		delete(pointers, v1ptr)
		return nil
	}
}

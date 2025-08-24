package sourceflag

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/krostar/cli"
	clicfg "github.com/krostar/cli/cfg"
)

// Source returns a SourceFunc that applies command-line flag values to a configuration struct.
//
// It uses reflection to find fields in the config struct that match the flag destinations
// and updates them with the flag values. It's designed to work with the context-based flag
// initialization system in the CLI package.
//
// The flagDest parameter is a pointer to the struct instance where flag values are stored.
// This is usually the same instance that contains the fields used as flag destinations in the
// command's Flags() and PersistentFlags() methods.
//
// The returned SourceFunc can be used in a BeforeCommandExecutionHook to load configuration
// from command-line flags. It should typically be the last source in the hook to ensure
// flag values take precedence over other configuration sources.
//
// Example:
//
//	type Config struct {
//		LogLevel string
//		Port     int
//	}
//
//	type MyCommand struct {
//		config Config
//	}
//
//	func (c *MyCommand) Flags() []cli.Flag {
//		return []cli.Flag{
//			cli.NewBuiltinFlag("log-level", "l", &c.config.LogLevel, "Set logging level"),
//			cli.NewBuiltinFlag("port", "p", &c.config.Port, "Set server port"),
//		}
//	}
//
//	func (c *MyCommand) Hook() *cli.Hook {
//		return &cli.Hook{
//			BeforeCommandExecution: clicfg.BeforeCommandExecutionHook(
//				&c.config,
//				// Apply flag values last for highest precedence
//				sourceflag.Source[Config](&c.config),
//			),
//		}
//	}
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
			return errors.New("some values where not found, make sure flag values all points to config")
		}

		return nil
	}
}

// recursivelyWalkThroughReflectValue recursively traverses two reflect.Values (v1 and v2)
// and updates fields in v2 with values from v1 based on the pointers in the map.
//
// This is an internal helper function that handles the complex task of walking through
// struct fields, following pointers, and copying values from flag destinations to the
// corresponding fields in the target configuration.
//
// Parameters:
//   - pointers: A map of memory addresses (as uintptr) for values set by flags.
//     When a matching field is found and processed, its address is removed from this map.
//   - v1: The reflect.Value of the flag destination struct containing flag values
//   - v2: The reflect.Value of the config struct where values should be copied to
//
// The function handles three main cases:
//  1. Pointers: It checks if the pointer itself is in the map, and if not, follows it
//  2. Structs: It recursively processes each field in the struct
//  3. Primitive types: It checks if the value's address is in the map, and if so, copies it
//
// It returns an error if any occurred during processing, such as an inability to address a value.
//
// At the end of successful processing, the pointers map should be empty, indicating that all
// flag values were successfully transferred to the config struct.
func recursivelyWalkThroughReflectValue(pointers map[uintptr]struct{}, v1, v2 reflect.Value, applyWOs ...func() error) error {
	if len(pointers) == 0 {
		return nil
	}

	switch k := v1.Type().Kind(); k {
	case reflect.Pointer:
		if v1.IsNil() {
			return nil
		}

		v1ptr := uintptr(v1.Addr().UnsafePointer())
		if _, ok := pointers[v1ptr]; !ok {
			var applyWO func() error

			v2, applyWO = ensurePointerInitialized(v2)
			applyWOs = append(applyWOs, applyWO)

			return recursivelyWalkThroughReflectValue(pointers, v1.Elem(), v2.Elem(), applyWOs...)
		}

		v2.Set(v1)
		delete(pointers, v1ptr)

		if err := applyAllWritingOperations(applyWOs); err != nil {
			return fmt.Errorf("unable to apply all writing operations: %v", err)
		}

		return nil

	case reflect.Struct:
		var errs []error
		for i := range v1.NumField() {
			errs = append(errs, recursivelyWalkThroughReflectValue(pointers, v1.Field(i), v2.Field(i), applyWOs...))
		}

		return errors.Join(errs...)

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

		if err := applyAllWritingOperations(applyWOs); err != nil {
			return fmt.Errorf("unable to apply all writing operations: %v", err)
		}

		return nil
	}
}

func ensurePointerInitialized(v reflect.Value) (reflect.Value, func() error) {
	if !v.IsNil() {
		return v, func() error { return nil }
	}

	oldV := v
	newV := reflect.New(v.Type().Elem())

	return newV, func() error {
		return reflectSetValue(oldV, newV)
	}
}

func reflectSetValue(dst, src reflect.Value) error {
	if !dst.IsValid() || !src.IsValid() {
		return fmt.Errorf("invalid value: dst.IsValid=%v src.IsValid=%v", dst.IsValid(), src.IsValid())
	}

	if dst.Type() != src.Type() {
		return fmt.Errorf("type mismatch: %v != %v", dst.Type(), src.Type())
	}

	if dst.CanSet() {
		dst.Set(src)
		return nil
	}

	if !dst.CanAddr() {
		return fmt.Errorf("dst is not addressable; type=%v", dst.Type())
	}

	dst = reflect.NewAt(dst.Type(), unsafe.Pointer(dst.UnsafeAddr())).Elem()
	dst.Set(src)

	return nil
}

func applyAllWritingOperations(wo []func() error) error {
	for i := range wo {
		if err := wo[len(wo)-1-i](); err != nil {
			return err
		}
	}

	return nil
}

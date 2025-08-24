package sourceenv

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/multierr"

	clicfg "github.com/krostar/cli/cfg"
)

// Source returns a SourceFunc that updates a config from environment variables.
// It uses reflection to traverse the config struct and sets fields based on environment variables.
// The environment variable names are derived from the struct field names, converted to uppercase and
// with nested fields separated by underscores. A struct field can specify additional environment
// variable names using the `env` tag, comma separated.
func Source[T any](envPrefix string) clicfg.SourceFunc[T] {
	return func(_ context.Context, cfg *T) error {
		_, err := recursivelyWalkThroughReflectValue(os.LookupEnv, reflect.ValueOf(cfg).Elem(), envPrefix, nil)
		return err
	}
}

// recursivelyWalkThroughReflectValue recursively traverses a reflect.Value and sets fields from environment variables.
//
//	lookupEnv is a function to lookup environment variables.
//	v is the current reflect.Value being processed.
//	envPrefix is the prefix for the environment variable names.
//	additionalEnvsToLookup are additional environment variable names specified in the `env` tag.
//
// It returns a boolean indicating if at least one environment variable was found and an error if any occurred.
func recursivelyWalkThroughReflectValue(lookupEnv func(string) (string, bool), v reflect.Value, envPrefix string, additionalEnvsToLookup []string) (bool, error) {
	t := v.Type()

	switch t.Kind() {
	case reflect.Pointer: // if it's a pointer, dereference it and continue recursively
		if !v.IsNil() {
			return recursivelyWalkThroughReflectValue(lookupEnv, v.Elem(), envPrefix, additionalEnvsToLookup)
		}

		// if the pointer is nil, create a new value of the underlying type
		newV := reflect.New(v.Type().Elem())
		// recursively process the new value
		atLeastOneEnvFound, err := recursivelyWalkThroughReflectValue(lookupEnv, newV.Elem(), envPrefix, additionalEnvsToLookup)
		// if at least one environment variable was found for the nested struct, set the pointer
		if atLeastOneEnvFound {
			v.Set(newV)
		}
		return atLeastOneEnvFound, err

	case reflect.Struct: // if it's a struct, iterate over its fields
		var (
			errs            []error
			atLeastOneFound bool
		)
		for i := range v.NumField() {
			tfield := t.Field(i)
			tag := tfield.Tag.Get("env")

			embededStruct := tfield.Anonymous && tfield.Type.Kind() == reflect.Struct
			unexported := tfield.PkgPath != ""
			skipped := tag == "-"

			// skip unexported fields that are not embedded structs, and fields with `env:"-"`
			if skipped || (unexported && !embededStruct) {
				continue
			}

			newEnvPrefix := envPrefix + "_" + strings.ToUpper(tfield.Name)
			if embededStruct && tag == "^" {
				newEnvPrefix = envPrefix
			}

			// recursively process each field, constructing the environment variable name
			envFound, err := recursivelyWalkThroughReflectValue(lookupEnv, v.Field(i), newEnvPrefix, strings.Split(tag, ","))
			if envFound {
				atLeastOneFound = true
			}
			errs = append(errs, err)
		}
		return atLeastOneFound, multierr.Combine(errs...)

	default: // for primitive types, try to find the corresponding environment variable
		var rawEnv string
		for _, envToLookup := range append(additionalEnvsToLookup, envPrefix) {
			envToLookup = strings.TrimSpace(envToLookup)
			if envToLookup != "" {
				if env, isset := lookupEnv(SanitizeName(envToLookup)); isset {
					rawEnv = env
					break
				}
			}
		}
		// no environment variable is found, return
		if rawEnv == "" {
			return false, nil
		}

		// convert the environment variable value to the field's type and set it
		switch k := t.Kind(); k {
		case reflect.Bool:
			vv, err := strconv.ParseBool(rawEnv)
			v.SetBool(vv)
			return true, err
		case reflect.Int:
			vv, err := strconv.ParseInt(rawEnv, 10, 0)
			v.SetInt(vv)
			return true, err
		case reflect.Int8:
			vv, err := strconv.ParseInt(rawEnv, 10, 8)
			v.SetInt(vv)
			return true, err
		case reflect.Int16:
			vv, err := strconv.ParseInt(rawEnv, 10, 16)
			v.SetInt(vv)
			return true, err
		case reflect.Int32:
			vv, err := strconv.ParseInt(rawEnv, 10, 32)
			v.SetInt(vv)
			return true, err
		case reflect.Int64:
			vv, err := strconv.ParseInt(rawEnv, 10, 64)
			v.SetInt(vv)
			return true, err
		case reflect.Uint:
			vv, err := strconv.ParseUint(rawEnv, 10, 0)
			v.SetUint(vv)
			return true, err
		case reflect.Uint8:
			vv, err := strconv.ParseUint(rawEnv, 10, 8)
			v.SetUint(vv)
			return true, err
		case reflect.Uint16:
			vv, err := strconv.ParseUint(rawEnv, 10, 16)
			v.SetUint(vv)
			return true, err
		case reflect.Uint32:
			vv, err := strconv.ParseUint(rawEnv, 10, 32)
			v.SetUint(vv)
			return true, err
		case reflect.Uint64:
			vv, err := strconv.ParseUint(rawEnv, 10, 64)
			v.SetUint(vv)
			return true, err
		case reflect.Float32:
			vv, err := strconv.ParseFloat(rawEnv, 32)
			v.SetFloat(vv)
			return true, err
		case reflect.Float64:
			vv, err := strconv.ParseFloat(rawEnv, 64)
			v.SetFloat(vv)
			return true, err
		case reflect.Complex64:
			vv, err := strconv.ParseComplex(rawEnv, 64)
			v.SetComplex(vv)
			return true, err
		case reflect.Complex128:
			vv, err := strconv.ParseComplex(rawEnv, 128)
			v.SetComplex(vv)
			return true, err
		case reflect.String:
			v.SetString(rawEnv)
			return true, nil
		default:
			return true, fmt.Errorf("unhandled type %s", k)
		}
	}
}

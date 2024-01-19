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

// Source updates config with environment variable value.
// A structure that can be set using A.B.C = 42 can be set
// through the environment variable named A_B_C = "42".
// Additional env values can be used with the `env` tag.
func Source[T any](envPrefix string) clicfg.SourceFunc[T] {
	return func(_ context.Context, cfg *T) error {
		_, err := recursivelyWalkThroughReflectValue(os.LookupEnv, reflect.ValueOf(cfg).Elem(), envPrefix, nil)
		return err
	}
}

func recursivelyWalkThroughReflectValue(lookupEnv func(string) (string, bool), v reflect.Value, envPrefix string, additionalEnvsToLookup []string) (bool, error) {
	t := v.Type()

	switch t.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			return recursivelyWalkThroughReflectValue(lookupEnv, v.Elem(), envPrefix, additionalEnvsToLookup)
		}

		newV := reflect.New(v.Type().Elem())
		atLeastOneEnvFound, err := recursivelyWalkThroughReflectValue(lookupEnv, newV.Elem(), envPrefix, additionalEnvsToLookup)
		if atLeastOneEnvFound {
			v.Set(newV)
		}
		return atLeastOneEnvFound, err

	case reflect.Struct:
		var (
			errs            []error
			atLeastOneFound bool
		)
		for i := 0; i < v.NumField(); i++ {
			tfield := t.Field(i)
			tag := tfield.Tag.Get("env")

			if tfield.PkgPath != "" || tag == "-" {
				continue
			}

			envFound, err := recursivelyWalkThroughReflectValue(lookupEnv, v.Field(i), envPrefix+"_"+strings.ToUpper(tfield.Name), strings.Split(tag, ","))
			if envFound {
				atLeastOneFound = true
			}
			errs = append(errs, err)
		}
		return atLeastOneFound, multierr.Combine(errs...)

	default:
		var rawEnv string
		for _, envToLookup := range append(additionalEnvsToLookup, envPrefix) {
			envToLookup = strings.TrimSpace(envToLookup)
			if envToLookup != "" {
				if env, isset := lookupEnv(envToLookup); isset {
					rawEnv = env
					break
				}
			}
		}
		if rawEnv == "" {
			return false, nil
		}

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

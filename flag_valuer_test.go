package cli

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func Test_NewFlagValuer(t *testing.T) {
	t.Run("non nil underlying value", func(t *testing.T) {
		var d time.Duration

		valuer := NewFlagValuer(&d, time.ParseDuration, func(d time.Duration) string { return d.String() })

		test.Assert(t, valuer.Destination() == &d)
		test.Assert(t, valuer.TypeRepr() == "time.Duration")
		test.Assert(t, !valuer.IsSet())

		err := valuer.FromString("abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid duration"))

		test.Assert(t, !valuer.IsSet())
		test.Assert(t, valuer.FromString("4s") == nil)
		test.Assert(t, valuer.IsSet())
		test.Assert(t, d == 4*time.Second)

		value := valuer.String()
		test.Assert(t, value == "4s")
	})

	t.Run("nil parameters", func(t *testing.T) {
		var d time.Duration

		test.Assert(check.Panics(t, func() {
			NewFlagValuer(nil, time.ParseDuration, func(d time.Duration) string { return d.String() })
		}, func(reason any) error {
			if strings.Contains(reason.(string), "destination is nil") {
				return nil
			}

			return errors.New("expected different panic reason")
		}))

		test.Assert(check.Panics(t, func() {
			NewFlagValuer(&d, nil, func(d time.Duration) string { return d.String() })
		}, func(reason any) error {
			if strings.Contains(reason.(string), "parse is nil") {
				return nil
			}

			return errors.New("expected different panic reason")
		}))

		test.Assert(check.Panics(t, func() {
			NewFlagValuer(&d, time.ParseDuration, nil)
		}, func(reason any) error {
			if strings.Contains(reason.(string), "toString is nil") {
				return nil
			}

			return errors.New("expected different panic reason")
		}))
	})
}

func Test_NewStringerFlagValuer(t *testing.T) {
	t.Run("non nil underlying value", func(t *testing.T) {
		var d time.Duration

		valuer := NewStringerFlagValuer(&d, time.ParseDuration)

		test.Assert(t, valuer.TypeRepr() == "time.Duration")
		test.Assert(t, !valuer.IsSet())

		err := valuer.FromString("abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid duration"))

		test.Assert(t, !valuer.IsSet())
		test.Assert(t, valuer.FromString("4s") == nil)
		test.Assert(t, valuer.IsSet())

		value := valuer.String()
		test.Assert(t, value == "4s")
	})

	t.Run("nil parameters", func(t *testing.T) {
		var d time.Duration

		test.Assert(check.Panics(t, func() {
			NewStringerFlagValuer(nil, time.ParseDuration)
		}, func(reason any) error {
			if strings.Contains(reason.(string), "destination is nil") {
				return nil
			}

			return errors.New("expected different panic reason")
		}))

		test.Assert(check.Panics(t, func() {
			NewStringerFlagValuer(&d, nil)
		}, func(reason any) error {
			if strings.Contains(reason.(string), "parse is nil") {
				return nil
			}

			return errors.New("expected different panic reason")
		}))
	})
}

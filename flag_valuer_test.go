package cli

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func Test_NewFlagValuer(t *testing.T) {
	t.Run("non nil underlying value", func(t *testing.T) {
		var d time.Duration

		valuer := NewFlagValuer(&d, time.ParseDuration, func(d time.Duration) string { return d.String() })

		assert.Check(t, valuer.TypeRepr() == "time.Duration")
		assert.Check(t, !valuer.IsSet())
		assert.ErrorContains(t, valuer.FromString("abc"), "invalid duration")
		assert.Check(t, !valuer.IsSet())
		assert.Check(t, valuer.FromString("4s"))
		assert.Check(t, valuer.IsSet())
		assert.Check(t, d == 4*time.Second)

		value := valuer.String()
		assert.Check(t, value == "4s")
	})

	t.Run("nil parameters", func(t *testing.T) {
		var d time.Duration

		assert.Check(t, func() (result cmp.Result) {
			defer func() {
				if reason := recover(); reason != nil {
					result = cmp.Equal(reason, "destination is nil")()
				}
			}()
			NewFlagValuer(nil, time.ParseDuration, func(d time.Duration) string { return d.String() })
			return cmp.ResultFailure("did not panic")
		})

		assert.Check(t, func() (result cmp.Result) {
			defer func() {
				if reason := recover(); reason != nil {
					result = cmp.Equal(reason, "parse is nil")()
				}
			}()
			NewFlagValuer(&d, nil, func(d time.Duration) string { return d.String() })
			return cmp.ResultFailure("did not panic")
		})

		assert.Check(t, func() (result cmp.Result) {
			defer func() {
				if reason := recover(); reason != nil {
					result = cmp.Equal(reason, "toString is nil")()
				}
			}()
			NewFlagValuer(&d, time.ParseDuration, nil)
			return cmp.ResultFailure("did not panic")
		})
	})
}

func Test_NewStringerFlagValuer(t *testing.T) {
	t.Run("non nil underlying value", func(t *testing.T) {
		var d time.Duration

		valuer := NewStringerFlagValuer(&d, time.ParseDuration)

		assert.Check(t, valuer.TypeRepr() == "time.Duration")
		assert.Check(t, !valuer.IsSet())
		assert.ErrorContains(t, valuer.FromString("abc"), "invalid duration")
		assert.Check(t, !valuer.IsSet())
		assert.Check(t, valuer.FromString("4s"))
		assert.Check(t, valuer.IsSet())

		value := valuer.String()
		assert.Check(t, value == "4s")
	})

	t.Run("nil parameters", func(t *testing.T) {
		var d time.Duration

		assert.Check(t, func() (result cmp.Result) {
			defer func() {
				if reason := recover(); reason != nil {
					result = cmp.Equal(reason, "destination is nil")()
				}
			}()
			NewStringerFlagValuer(nil, time.ParseDuration)
			return cmp.ResultFailure("did not panic")
		})

		assert.Check(t, func() (result cmp.Result) {
			defer func() {
				if reason := recover(); reason != nil {
					result = cmp.Equal(reason, "parse is nil")()
				}
			}()
			NewStringerFlagValuer(&d, nil)
			return cmp.ResultFailure("did not panic")
		})
	})
}

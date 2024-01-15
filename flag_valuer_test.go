package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewFlagValuer(t *testing.T) {
	t.Run("non nil underlying value", func(t *testing.T) {
		var d time.Duration

		valuer := NewFlagValuer(&d, time.ParseDuration, func(d time.Duration) string { return d.String() })

		assert.Equal(t, "time.Duration", valuer.TypeRepr())

		require.Error(t, valuer.FromString("abc"))
		require.NoError(t, valuer.FromString("4s"))
		assert.Equal(t, time.Second*4, d)

		value := valuer.String()
		assert.Equal(t, "4s", value)
	})

	t.Run("nil parameters", func(t *testing.T) {
		var d time.Duration

		assert.PanicsWithValue(t, "destination is nil", func() {
			NewFlagValuer(nil, time.ParseDuration, func(d time.Duration) string { return d.String() })
		})

		assert.PanicsWithValue(t, "parse is nil", func() {
			NewFlagValuer(&d, nil, func(d time.Duration) string { return d.String() })
		})

		assert.PanicsWithValue(t, "toString is nil", func() {
			NewFlagValuer(&d, time.ParseDuration, nil)
		})
	})
}

func Test_NewStringerFlagValuer(t *testing.T) {
	t.Run("non nil underlying value", func(t *testing.T) {
		var d time.Duration

		valuer := NewStringerFlagValuer(&d, time.ParseDuration)

		assert.Equal(t, "time.Duration", valuer.TypeRepr())

		require.Error(t, valuer.FromString("abc"))
		require.NoError(t, valuer.FromString("4s"))

		value := valuer.String()
		assert.Equal(t, "4s", value)
	})

	t.Run("nil parameters", func(t *testing.T) {
		var d time.Duration

		assert.PanicsWithValue(t, "destination is nil", func() {
			NewStringerFlagValuer(nil, time.ParseDuration)
		})

		assert.PanicsWithValue(t, "parse is nil", func() {
			NewStringerFlagValuer(&d, nil)
		})
	})
}

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFlag(t *testing.T) {
	var dest int
	nonNilValuer := &flagValuer[int]{value: &dest}

	t.Run("ok", func(t *testing.T) {
		assert.NotNil(t, NewFlag("long", "", nonNilValuer, "foo"))
		assert.NotNil(t, NewFlag("", "s", nonNilValuer, "foo"))
		assert.NotNil(t, NewFlag("long", "s", nonNilValuer, "foo"))
		assert.NotNil(t, NewFlag("long", "s", nonNilValuer, ""))
	})

	t.Run("wrong setup", func(t *testing.T) {
		t.Run("short and long names are unset", func(t *testing.T) {
			assert.PanicsWithValue(t, "longName and/or shortName must be non-empty", func() {
				NewFlag("", "", nonNilValuer, "")
			})
		})

		t.Run("short is set with more than one character", func(t *testing.T) {
			assert.PanicsWithValue(t, "shortName must be one character long", func() {
				NewFlag("long", "notSoShort", nonNilValuer, "")
			})
		})

		t.Run("nil valuer", func(t *testing.T) {
			assert.PanicsWithValue(t, "a non-nil valuer is required", func() {
				NewFlag("long", "", nil, "")
			})
		})
	})
}

func Test_flagValue_LongName(t *testing.T) {
	assert.Equal(t, "long", flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.LongName())
}

func Test_flagValue_ShortName(t *testing.T) {
	assert.Equal(t, "short", flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.ShortName())
}

func Test_flagValue_Description(t *testing.T) {
	assert.Equal(t, "description", flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.Description())
}

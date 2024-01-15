package cli

import (
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func Test_NewFlag(t *testing.T) {
	var dest int
	nonNilValuer := &flagValuer[int]{value: &dest}

	t.Run("ok", func(t *testing.T) {
		assert.Check(t, NewFlag("long", "", nonNilValuer, "foo") != nil)
		assert.Check(t, NewFlag("", "s", nonNilValuer, "foo") != nil)
		assert.Check(t, NewFlag("long", "s", nonNilValuer, "foo") != nil)
		assert.Check(t, NewFlag("long", "s", nonNilValuer, "") != nil)
	})

	t.Run("wrong setup", func(t *testing.T) {
		t.Run("short and long names are unset", func(t *testing.T) {
			assert.Check(t, func() (result cmp.Result) {
				defer func() {
					if reason := recover(); reason != nil {
						result = cmp.Equal(reason, "longName and/or shortName must be non-empty")()
					}
				}()
				NewFlag("", "", nonNilValuer, "")
				return cmp.ResultFailure("did not panic")
			})
		})

		t.Run("short is set with more than one character", func(t *testing.T) {
			assert.Check(t, func() (result cmp.Result) {
				defer func() {
					if reason := recover(); reason != nil {
						result = cmp.Equal(reason, "shortName must be one character long")()
					}
				}()
				NewFlag("long", "notSoShort", nonNilValuer, "")
				return cmp.ResultFailure("did not panic")
			})
		})

		t.Run("nil valuer", func(t *testing.T) {
			assert.Check(t, func() (result cmp.Result) {
				defer func() {
					if reason := recover(); reason != nil {
						result = cmp.Equal(reason, "a non-nil valuer is required")()
					}
				}()
				NewFlag("long", "", nil, "")
				return cmp.ResultFailure("did not panic")
			})
		})
	})
}

func Test_flagValue_LongName(t *testing.T) {
	assert.Check(t, flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.LongName() == "long")
}

func Test_flagValue_ShortName(t *testing.T) {
	assert.Check(t, flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.ShortName() == "short")
}

func Test_flagValue_Description(t *testing.T) {
	assert.Check(t, flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.Description() == "description")
}

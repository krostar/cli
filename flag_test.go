package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func Test_NewFlag(t *testing.T) {
	var dest int
	nonNilValuer := &flagValuer[int]{value: &dest}

	t.Run("ok", func(t *testing.T) {
		test.Assert(t, NewFlag("long", "", nonNilValuer, "foo") != nil)
		test.Assert(t, NewFlag("", "s", nonNilValuer, "foo") != nil)
		test.Assert(t, NewFlag("long", "s", nonNilValuer, "foo") != nil)
		test.Assert(t, NewFlag("long", "s", nonNilValuer, "") != nil)
	})

	t.Run("wrong setup", func(t *testing.T) {
		t.Run("short and long names are unset", func(t *testing.T) {
			test.Assert(check.Panics(t, func() {
				NewFlag("", "", nonNilValuer, "")
			}, func(reason any) error {
				if strings.Contains(reason.(string), "longName and/or shortName must be non-empty") {
					return nil
				}
				return errors.New("expected different panic reason")
			}))
		})

		t.Run("short is set with more than one character", func(t *testing.T) {
			test.Assert(check.Panics(t, func() {
				NewFlag("long", "notSoShort", nonNilValuer, "")
			}, func(reason any) error {
				if strings.Contains(reason.(string), "shortName must be one character long") {
					return nil
				}
				return errors.New("expected different panic reason")
			}))
		})

		t.Run("nil valuer", func(t *testing.T) {
			test.Assert(check.Panics(t, func() {
				NewFlag("long", "", nil, "")
			}, func(reason any) error {
				if strings.Contains(reason.(string), "a non-nil valuer is required") {
					return nil
				}
				return errors.New("expected different panic reason")
			}))
		})
	})
}

func Test_flagValue_LongName(t *testing.T) {
	test.Assert(t, flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.LongName() == "long")
}

func Test_flagValue_ShortName(t *testing.T) {
	test.Assert(t, flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.ShortName() == "short")
}

func Test_flagValue_Description(t *testing.T) {
	test.Assert(t, flagValue{
		longName:    "long",
		shortName:   "short",
		description: "description",
	}.Description() == "description")
}

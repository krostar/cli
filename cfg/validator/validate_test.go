package clicfgvalidator

import (
	"errors"
	"testing"

	"github.com/krostar/test"
)

func Test_Validate(t *testing.T) {
	t.Run("cfg without Validate method", func(t *testing.T) {
		cfg := configWithoutValidate{A: "foo"}
		test.Assert(t, Validate(&cfg) == nil)
	})

	t.Run("cfg with Validate returning nil", func(t *testing.T) {
		cfg := configWithValidate{A: "foo"}
		test.Assert(t, Validate(&cfg) == nil)
	})

	t.Run("cfg with Validate returning error", func(t *testing.T) {
		var cfg configWithValidate
		test.Assert(t, errors.Is(Validate(&cfg), errInvalid))
	})

	t.Run("nested field Validate is called when parent has no Validate", func(t *testing.T) {
		var cfg configWithNestedValidate
		test.Assert(t, errors.Is(Validate(&cfg), errInvalid))
	})

	t.Run("nested field Validate is not called when parent implements Validate", func(t *testing.T) {
		var cfg configOwnsValidation // Sub.B is empty but parent's Validate ignores it
		test.Assert(t, Validate(&cfg) == nil)
	})

	t.Run("nil pointer field is skipped", func(t *testing.T) {
		var cfg configWithNilPtr
		test.Assert(t, Validate(&cfg) == nil)
	})

	t.Run("non-nil pointer field Validate is called", func(t *testing.T) {
		cfg := configWithNilPtr{Sub: &subConfigWithValidate{}}
		test.Assert(t, errors.Is(Validate(&cfg), errInvalid))
	})
}

var errInvalid = errors.New("invalid")

type configWithoutValidate struct {
	A string
}

type configWithValidate struct {
	A string
}

func (cfg *configWithValidate) Validate() error {
	if cfg.A == "" {
		return errInvalid
	}
	return nil
}

type configWithNestedValidate struct {
	Sub subConfigWithValidate
}

type subConfigWithValidate struct {
	B string
}

func (cfg *subConfigWithValidate) Validate() error {
	if cfg.B == "" {
		return errInvalid
	}
	return nil
}

// configOwnsValidation implements Validate itself, so its Sub field must NOT be visited.
type configOwnsValidation struct {
	Sub subConfigWithValidate
}

func (cfg *configOwnsValidation) Validate() error {
	return nil // deliberately ignores Sub.B being empty
}

type configWithNilPtr struct {
	Sub *subConfigWithValidate
}

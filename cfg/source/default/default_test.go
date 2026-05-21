package sourcedefault

import (
	"testing"

	"github.com/krostar/test"
)

func Test_Source(t *testing.T) {
	t.Run("cfg without SetDefault method", func(t *testing.T) {
		var cfg configWithoutDefault

		test.Require(t, Source[configWithoutDefault]()(test.Context(t), &cfg) == nil)
		test.Assert(t, cfg == configWithoutDefault{})
	})

	t.Run("cfg with SetDefault method", func(t *testing.T) {
		var cfg configWithDefault

		test.Require(t, Source[configWithDefault]()(test.Context(t), &cfg) == nil)
		test.Assert(t, cfg == configWithDefault{A: "foo"})
	})

	t.Run("cfg with nested struct field implementing SetDefault", func(t *testing.T) {
		var cfg configWithNestedDefault

		test.Require(t, Source[configWithNestedDefault]()(test.Context(t), &cfg) == nil)
		test.Assert(t, cfg.Sub == subConfigWithDefault{B: "bar"})
	})

	t.Run("cfg with nil pointer field implementing SetDefault gets initialized", func(t *testing.T) {
		var cfg configWithNestedDefault

		test.Require(t, Source[configWithNestedDefault]()(test.Context(t), &cfg) == nil)
		test.Require(t, cfg.SubPtr != nil)
		test.Assert(t, *cfg.SubPtr == subConfigWithDefault{B: "bar"})
	})

	t.Run("cfg with non-nil pointer field implementing SetDefault", func(t *testing.T) {
		cfg := configWithNestedDefault{SubPtr: &subConfigWithDefault{}}

		test.Require(t, Source[configWithNestedDefault]()(test.Context(t), &cfg) == nil)
		test.Assert(t, *cfg.SubPtr == subConfigWithDefault{B: "bar"})
	})

	t.Run("cfg with deeply nested field implementing SetDefault", func(t *testing.T) {
		var cfg configWithDeepDefault

		test.Require(t, Source[configWithDeepDefault]()(test.Context(t), &cfg) == nil)
		test.Assert(t, cfg.Sub.Deep == subConfigWithDefault{B: "bar"})
	})

	t.Run("nil pointer to non-setter type is initialized when a nested field has SetDefault", func(t *testing.T) {
		var cfg configWithNilPtrToNonSetter

		test.Require(t, Source[configWithNilPtrToNonSetter]()(test.Context(t), &cfg) == nil)
		test.Require(t, cfg.Sub != nil)
		test.Assert(t, cfg.Sub.Inner == subConfigWithDefault{B: "bar"})
	})

	t.Run("nil pointer to non-setter type stays nil when no nested field has SetDefault", func(t *testing.T) {
		var cfg struct{ Sub *configWithoutDefault }

		test.Require(t, Source[struct{ Sub *configWithoutDefault }]()(test.Context(t), &cfg) == nil)
		test.Assert(t, cfg.Sub == nil)
	})

	t.Run("parent SetDefault runs after field SetDefault and can override it", func(t *testing.T) {
		var cfg configOverridesField

		test.Require(t, Source[configOverridesField]()(test.Context(t), &cfg) == nil)
		// parent overrides B, child's C is untouched
		test.Assert(t, cfg.Sub == subConfigOverridable{B: "parent", C: "child"})
	})
}

type configWithoutDefault struct {
	A string
}

type configWithDefault struct {
	A string
}

func (cfg *configWithDefault) SetDefault() {
	cfg.A = "foo"
}

type configWithNestedDefault struct {
	Sub    subConfigWithDefault
	SubPtr *subConfigWithDefault
}

type subConfigWithDefault struct {
	B string
}

func (cfg *subConfigWithDefault) SetDefault() {
	cfg.B = "bar"
}

type configWithDeepDefault struct {
	Sub subConfigWithDeepDefault
}

type subConfigWithDeepDefault struct {
	Deep subConfigWithDefault
}

type configWithNilPtrToNonSetter struct {
	Sub *subWithNestedDefault
}

type subWithNestedDefault struct {
	Inner subConfigWithDefault
}

// configOverridesField has both a root SetDefault and a field with its own SetDefault.
// The root's SetDefault runs after the field's (bottom-up), so it wins on B but leaves C alone.
type configOverridesField struct {
	Sub subConfigOverridable
}

type subConfigOverridable struct {
	B string
	C string
}

func (cfg *subConfigOverridable) SetDefault() {
	cfg.B = "child"
	cfg.C = "child"
}

func (cfg *configOverridesField) SetDefault() {
	cfg.Sub.B = "parent" // intentionally overrides child's B
}

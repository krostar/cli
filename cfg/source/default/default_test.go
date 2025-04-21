package sourcedefault

import (
	"testing"

	"github.com/krostar/test"
)

type configWithoutDefault struct {
	A string
}

type configWithDefault struct {
	A string
}

func (cfg *configWithDefault) SetDefault() {
	cfg.A = "foo"
}

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
}

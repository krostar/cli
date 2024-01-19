package sourcedefault

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
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

		assert.NilError(t, Source[configWithoutDefault]()(context.Background(), &cfg))
		assert.DeepEqual(t, cfg, configWithoutDefault{})
	})

	t.Run("cfg with SetDefault method", func(t *testing.T) {
		var cfg configWithDefault

		assert.NilError(t, Source[configWithDefault]()(context.Background(), &cfg))
		assert.DeepEqual(t, cfg, configWithDefault{A: "foo"})
	})
}

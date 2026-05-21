package clicfg

import (
	"context"
	"errors"
	"testing"

	"github.com/krostar/test"
)

func Test_Apply(t *testing.T) {
	type config struct{ A string }

	t.Run("sources are applied in order", func(t *testing.T) {
		var cfg config

		test.Require(t, Apply(
			test.Context(t), &cfg,
			func(_ context.Context, cfg *config) error { cfg.A += "1"; return nil },
			func(_ context.Context, cfg *config) error { cfg.A += "2"; return nil },
			func(_ context.Context, cfg *config) error { cfg.A += "3"; return nil },
		) == nil)

		test.Assert(t, cfg.A == "123")
	})

	t.Run("error is returned and wrapped with source index", func(t *testing.T) {
		var cfg config

		expectedErr := errors.New("boom")

		err := Apply(
			test.Context(t), &cfg,
			func(context.Context, *config) error { return nil },
			func(context.Context, *config) error { return expectedErr },
		)

		test.Assert(t, errors.Is(err, expectedErr))
	})
}

func Test_BeforeCommandExecutionHook(t *testing.T) {
	type config struct{ A string }

	var cfg config

	test.Require(t, BeforeCommandExecutionHook(
		&cfg,
		func(_ context.Context, cfg *config) error { cfg.A = "ok"; return nil },
	)(test.Context(t)) == nil)

	test.Assert(t, cfg.A == "ok")
}

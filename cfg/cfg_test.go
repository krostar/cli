package clicfg

import (
	"context"
	"errors"
	"testing"

	"github.com/krostar/test"
)

func Test_BeforeCommandExecutionHook(t *testing.T) {
	type config struct {
		A string
	}
	var cfg config

	t.Run("ok", func(t *testing.T) {
		test.Require(t, BeforeCommandExecutionHook(&cfg,
			func(_ context.Context, cfg *config) error {
				cfg.A += "1"
				return nil
			},
			func(_ context.Context, cfg *config) error {
				cfg.A += "2"
				return nil
			},
			func(_ context.Context, cfg *config) error {
				cfg.A += "3"
				return nil
			},
		)(test.Context(t)) == nil)

		test.Assert(t, cfg.A == "123")
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("boom")

		test.Assert(t, errors.Is(BeforeCommandExecutionHook(&cfg,
			func(context.Context, *config) error {
				return expectedErr
			},
		)(test.Context(t)), expectedErr))
	})
}

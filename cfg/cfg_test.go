package clicfg

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_BeforeCommandExecutionHook(t *testing.T) {
	type config struct {
		A string
	}
	var cfg config

	t.Run("ok", func(t *testing.T) {
		assert.NilError(t, BeforeCommandExecutionHook(&cfg,
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
		)(context.Background()))

		assert.Check(t, cfg.A == "123")
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("boom")

		assert.ErrorIs(t, BeforeCommandExecutionHook(&cfg,
			func(context.Context, *config) error {
				return expectedErr
			},
		)(context.Background()), expectedErr)
	})
}

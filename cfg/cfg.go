package clicfg

import (
	"context"
	"fmt"

	"github.com/krostar/cli"
)

// SourceFunc defines the function signature to apply a config source to the provided config.
type SourceFunc[T any] func(ctx context.Context, cfg *T) error

// BeforeCommandExecutionHook replaces dest with provided config sources.
func BeforeCommandExecutionHook[T any](dest *T, source SourceFunc[T], sources ...SourceFunc[T]) cli.HookFunc {
	sources = append([]SourceFunc[T]{source}, sources...)

	return func(ctx context.Context) error {
		cfg := new(T)

		for i, source := range sources {
			if err := source(ctx, cfg); err != nil {
				return fmt.Errorf("unable to apply config source[%d]: %w", i, err)
			}
		}

		*dest = *cfg
		return nil
	}
}

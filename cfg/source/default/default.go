package sourcedefault

import (
	"context"

	clicfg "github.com/krostar/cli/cfg"
)

// Source calls the SetDefault() method from provided Config, if it exists.
func Source[T any]() clicfg.SourceFunc[T] {
	return func(_ context.Context, cfg *T) error {
		if cfgDefault, ok := (any(cfg)).(interface {
			SetDefault()
		}); ok {
			cfgDefault.SetDefault()
		}
		return nil
	}
}

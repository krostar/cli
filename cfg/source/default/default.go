package sourcedefault

import (
	"context"

	clicfg "github.com/krostar/cli/cfg"
)

// Source returns a SourceFunc that sets default values for a config.
// It checks if the config implements a SetDefault() method and calls it if available.
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

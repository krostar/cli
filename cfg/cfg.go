// Package clicfg provides a flexible configuration system for CLI applications.
// It allows loading configuration from multiple sources (default values, files,
// environment variables, command-line flags) with a clear precedence order.
package clicfg

import (
	"context"
	"fmt"

	"github.com/krostar/cli"
)

// SourceFunc defines a function signature for applying a configuration source to a given config.
// A SourceFunc takes a context and a pointer to a config struct and populates the config
// from its specific source (e.g., environment variables, config file, etc.).
//
// The type parameter T represents the configuration struct type.
// Each source implementation knows how to populate fields of type T from its specific source.
type SourceFunc[T any] func(ctx context.Context, cfg *T) error

// BeforeCommandExecutionHook creates a hook that applies multiple configuration sources
// to a destination config. The sources are applied in the order they are provided,
// with later sources overriding values from earlier ones if they provide the same setting.
// Returns a cli.HookFunc that can be used as a BeforeCommandExecution hook.
//
// Example:
//
//	func (cmd *MyCommand) Hook() *cli.Hook {
//	    return &cli.Hook{
//	        BeforeCommandExecution: clicfg.BeforeCommandExecutionHook(
//	            &cmd.config,
//	            sourcedefault.Source[Config](),                 // 1. Default values (lowest priority)
//	            sourcefile.Source(getFilePath, decoder, true),  // 2. Config from file
//	            sourceenv.Source[Config]("APP"),                // 3. Environment variables
//	            sourceflag.Source[Config](cmd),                 // 4. Command-line flags (highest priority)
//	        ),
//	    }
//	}
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

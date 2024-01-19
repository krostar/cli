package sourcefile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	clicfg "github.com/krostar/cli/cfg"
)

// Source opens the file provided by getFilename() and calls unmarshaler.
// If allowNonExisting and the file does not exist, Source do nothing.
func Source[T any](getFilename func(cfg T) string, unmarshaler func(reader io.Reader, cfg *T) error, allowNonExisting bool) clicfg.SourceFunc[T] {
	return func(_ context.Context, cfg *T) error {
		file, err := os.Open(getFilename(*cfg))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) && allowNonExisting {
				return nil
			}
			return fmt.Errorf("unable to open config file: %w", err)
		}
		defer file.Close() //nolint:errcheck // standard library don't care about this error

		if err := unmarshaler(file, cfg); err != nil {
			return fmt.Errorf("unable to decode config: %w", err)
		}

		return nil
	}
}

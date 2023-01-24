package cli

import (
	"context"
	"os"
	"os/signal"
)

type (
	ctxKey      uint8
	ctxMetadata map[any]any
)

const (
	ctxKeyMetadata ctxKey = iota
	ctxKeyExitLogger
)

// NewContextWithMetadata wraps the provided context to create a global metadata store to the CLI.
// It helps to pass dependencies across commands tree.
func NewContextWithMetadata(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyMetadata, make(ctxMetadata))
}

// SetMetadataInContext associates a key to a value in the global CLI metadata store.
func SetMetadataInContext(ctx context.Context, key, value any) {
	if meta, ok := ctx.Value(ctxKeyMetadata).(ctxMetadata); ok {
		meta[key] = value
	}
}

// GetMetadataFromContext retrieves any value stores to the provided key, if any.
func GetMetadataFromContext(ctx context.Context, key any) any {
	if meta, ok := ctx.Value(ctxKeyMetadata).(ctxMetadata); ok {
		return meta[key]
	}
	return nil
}

// NewContextCancelableBySignal creates a new context that cancels itself when provided signals are triggered.
func NewContextCancelableBySignal(signals ...os.Signal) (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = NewContextWithMetadata(ctx)

	if len(signals) == 0 {
		return ctx, cancel
	}

	signalChan := make(chan os.Signal, 1)
	clean := func() {
		signal.Ignore(signals...)
		close(signalChan)
	}

	// catch some stop signals, and cancel the context if caught
	signal.Notify(signalChan, signals...)
	go func() {
		<-signalChan // block until a signal is received
		cancel()
	}()

	return ctx, clean
}

package cli

import (
	"context"
	"os"
	"os/signal"
)

type ctxKey uint8

const (
	ctxMetadataKey ctxKey = iota
	ctxExitLogger
)

// ContextWithMetadata wraps the provided context to add a global metadata store to the CLI.
func ContextWithMetadata(ctx context.Context) context.Context {
	metadata := make(map[interface{}]interface{})
	ctx = context.WithValue(ctx, ctxMetadataKey, metadata)
	return ctx
}

// SetMetadata associates a key to a value in the global CLI metadata store.
func SetMetadata(ctx context.Context, key interface{}, value interface{}) {
	if meta, ok := ctx.Value(ctxMetadataKey).(map[interface{}]interface{}); ok {
		meta[key] = value
	}
}

// GetMetadata retrieves any value stores to the provided key, if any.
func GetMetadata(ctx context.Context, key interface{}) interface{} {
	if meta, ok := ctx.Value(ctxMetadataKey).(map[interface{}]interface{}); ok {
		return meta[key]
	}
	return nil
}

// NewContextCancelableBySignal creates a new context that cancels itself when provided signals are triggered.
func NewContextCancelableBySignal(signals ...os.Signal) (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = ContextWithMetadata(ctx)

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

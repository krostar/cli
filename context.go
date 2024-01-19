package cli

import (
	"context"
	"os"
	"os/signal"
)

type (
	ctxKey     uint8
	ctxCommand struct {
		localFlags      []Flag
		persistentFlags []Flag
	}
	ctxMetadata map[any]any
)

const (
	ctxKeyCommand ctxKey = iota
	ctxKeyMetadata
	metadataKeyExitLogger
)

// NewCommandContext is called for each command to create a dedicated context.
func NewCommandContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyCommand, new(ctxCommand))
}

func getCommandFromContext(ctx context.Context) *ctxCommand {
	if cmd, ok := ctx.Value(ctxKeyCommand).(*ctxCommand); ok {
		return cmd
	}
	return nil
}

// SetInitializedFlagsInContext sets the provided initialized flags in the command context.
func SetInitializedFlagsInContext(ctx context.Context, localFlags, persistentFlags []Flag) {
	if cmd := getCommandFromContext(ctx); cmd != nil {
		cmd.localFlags = localFlags
		cmd.persistentFlags = persistentFlags
	}
}

// GetInitializedFlagsFromContext returns initialized command flags.
func GetInitializedFlagsFromContext(ctx context.Context) ([]Flag, []Flag) {
	if cmd, ok := ctx.Value(ctxKeyCommand).(*ctxCommand); ok {
		return cmd.localFlags, cmd.persistentFlags
	}
	return nil, nil
}

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
func NewContextCancelableBySignal(sig os.Signal, sigs ...os.Signal) (context.Context, func()) {
	signals := append([]os.Signal{sig}, sigs...)

	ctx, cancel := context.WithCancel(context.Background())
	ctx = NewContextWithMetadata(ctx)

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

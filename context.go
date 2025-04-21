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
	ctxKeyCommand ctxKey = iota + 1
	ctxKeyMetadata
	metadataKeyExitLogger
)

// NewCommandContext is called for each command to create a dedicated context.
// Warning: This does not make sens to use outside of cli mapper.
func NewCommandContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyCommand, new(ctxCommand))
}

func getCommandFromContext(ctx context.Context) *ctxCommand {
	if cmd, ok := ctx.Value(ctxKeyCommand).(*ctxCommand); ok {
		return cmd
	}
	return nil
}

// SetInitializedFlagsInContext stores the initialized local and persistent
// flags for a command in the context. This allows the flags to be accessed
// later, for example, by configuration sources.
// Warning: This function is exposed for cli mappers, you should not use it directly.
func SetInitializedFlagsInContext(ctx context.Context, localFlags, persistentFlags []Flag) {
	if cmd := getCommandFromContext(ctx); cmd != nil {
		cmd.localFlags = localFlags
		cmd.persistentFlags = persistentFlags
	}
}

// GetInitializedFlagsFromContext retrieves the initialized local and
// persistent flags for a command from the context. Returns nil slices if
// no flags are found.
// Warning: This function is exposed for sourcing flag values in configuration loader, you should not use it directly.
func GetInitializedFlagsFromContext(ctx context.Context) ([]Flag, []Flag) {
	if cmd, ok := ctx.Value(ctxKeyCommand).(*ctxCommand); ok {
		return cmd.localFlags, cmd.persistentFlags
	}
	return nil, nil
}

// NewContextWithMetadata creates a new context that includes a metadata
// store. This store can be used to pass arbitrary data between different
// parts of the CLI application. This function should be called at the
// beginning of the CLI application's execution.
func NewContextWithMetadata(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyMetadata, make(ctxMetadata))
}

// SetMetadataInContext associates a key-value pair in the global CLI
// metadata store. This allows for storing and retrieving data that needs
// to be accessible across the entire CLI application.
func SetMetadataInContext(ctx context.Context, key, value any) {
	if meta, ok := ctx.Value(ctxKeyMetadata).(ctxMetadata); ok {
		meta[key] = value
	}
}

// GetMetadataFromContext retrieves a value from the global CLI metadata
// store based on the provided key. Returns nil if the key is not found.
func GetMetadataFromContext(ctx context.Context, key any) any {
	if meta, ok := ctx.Value(ctxKeyMetadata).(ctxMetadata); ok {
		return meta[key]
	}
	return nil
}

// NewContextCancelableBySignal creates a new context that is automatically
// canceled when any of the provided signals are received. This is useful
// for gracefully shutting down the CLI application on interrupt signals
// (e.g., Ctrl+C).
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

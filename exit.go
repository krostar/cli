package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

// Exit exits the program and uses provided error to define program success or failure.
func Exit(ctx context.Context, err error, options ...ExitOption) {
	o := exitOptions{
		exitFunc:      os.Exit,
		getLoggerFunc: getExitLoggerFromMetadata,
	}

	for _, option := range options {
		option(&o)
	}

	var (
		msg    string
		status uint8
	)

	if err != nil {
		var errWithStatus ExitStatusError
		if errors.As(err, &errWithStatus) {
			status = errWithStatus.ExitStatus()
		} else {
			status = 255
		}

		msg = err.Error()
	}

	writer := o.getLoggerFunc(ctx)

	if msg != "" {
		if _, err := io.WriteString(writer, msg+"\n"); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("unable to write program exit message: %v", err)) //nolint:errcheck // we can't properly handle writing error on stderr
		}
	}

	if err := writer.Close(); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("unable to close writer: %v", err)) //nolint:errcheck // we can't properly handle writing error on stderr
	}

	o.exitFunc(int(status))
}

type exitOptions struct {
	exitFunc      func(int)
	getLoggerFunc func(context.Context) io.WriteCloser
}

// ExitOption defines the function signature to configure things upon exit.
type ExitOption func(*exitOptions)

// WithExitFunc defines a custom function to exit.
func WithExitFunc(exitFunc func(status int)) ExitOption {
	return func(o *exitOptions) {
		o.exitFunc = exitFunc
	}
}

// WithExitLoggerFunc defines a way to customize logger used in messages.
func WithExitLoggerFunc(getLoggerFunc func(context.Context) io.WriteCloser) ExitOption {
	return func(o *exitOptions) {
		o.getLoggerFunc = getLoggerFunc
	}
}

// SetExitLoggerInMetadata sets the logger used by the CLI to write the exit message if any, inside the metadata.
// By default, the Exit func tries to find the logger in the metadata.
func SetExitLoggerInMetadata(ctx context.Context, writer io.WriteCloser) {
	SetMetadataInContext(ctx, ctxKeyExitLogger, writer)
}

func getExitLoggerFromMetadata(ctx context.Context) io.WriteCloser {
	rawWriter := GetMetadataFromContext(ctx, ctxKeyExitLogger)
	if writer, ok := rawWriter.(io.WriteCloser); ok {
		return writer
	}
	return os.Stderr
}

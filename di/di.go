package clidi

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/dig"
	"go.uber.org/multierr"

	"github.com/krostar/cli"
)

type contextKey int8

const (
	contextKeyUnknown contextKey = iota
	contextKeyDIContainer
	contextKeyDIProvideErrors
)

type (
	// In alias dig.In to have only one import for DI in most case.
	In = dig.In
	// Out alias dig.Out to have only one import for DI in most case.
	Out = dig.Out
)

// InitializeContainer creates the DI container and sets it in the context.
func InitializeContainer(ctx context.Context, opts ...dig.Option) {
	cli.SetMetadataInContext(ctx, contextKeyDIContainer, dig.New(opts...))
}

// AddProvider provides a function f to the list of DI providers set in the DI container.
func AddProvider(ctx context.Context, f any, opts ...dig.ProvideOption) {
	container, exists := cli.GetMetadataFromContext(ctx, contextKeyDIContainer).(*dig.Container)
	if !exists {
		return
	}

	if err := container.Provide(f, opts...); err != nil {
		errs, _ := cli.GetMetadataFromContext(ctx, contextKeyDIProvideErrors).([]error)
		cli.SetMetadataInContext(ctx, contextKeyDIProvideErrors, append(errs, err))
	}
}

// Invoke uses previously provided functions (through AddProvider) to initialize function f.
func Invoke(ctx context.Context, f any, opts ...dig.InvokeOption) error {
	container, exists := cli.GetMetadataFromContext(ctx, contextKeyDIContainer).(*dig.Container)
	if !exists {
		return errors.New("container is unset in the context")
	}

	if errs, _ := cli.GetMetadataFromContext(ctx, contextKeyDIProvideErrors).([]error); len(errs) > 0 {
		return fmt.Errorf("provider error: %v", multierr.Combine(errs...))
	}

	if err := container.Invoke(f, opts...); err != nil {
		return fmt.Errorf("invoker error: %v", dig.RootCause(err))
	}

	return nil
}

package clidi

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/dig"

	"github.com/krostar/cli"
)

type contextKey int8

const (
	contextKeyDIContainer contextKey = iota + 1
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

// AddProvider adds a constructor function (f) to the dig.Container stored in
// the context. Any errors encountered during provider registration are
// collected and can be retrieved later using Invoke.
func AddProvider(ctx context.Context, f any, opts ...dig.ProvideOption) {
	container, exists := cli.GetMetadataFromContext(ctx, contextKeyDIContainer).(*dig.Container)
	if !exists {
		return
	}

	if err := container.Provide(f, opts...); err != nil {
		errs, _ := cli.GetMetadataFromContext(ctx, contextKeyDIProvideErrors).([]error) //nolint:revive,errcheck // unchecked-type-assertion: we know this type for sure
		cli.SetMetadataInContext(ctx, contextKeyDIProvideErrors, append(errs, err))
	}
}

// Invoke executes a function (f) using the dig.Container stored in the
// context. It first checks for any errors that occurred during provider
// registration and returns them if found. If the container is not found
// in the context, it returns an error. Otherwise, it invokes the function
// using the container, handling dependency injection.
func Invoke(ctx context.Context, f any, opts ...dig.InvokeOption) error {
	container, exists := cli.GetMetadataFromContext(ctx, contextKeyDIContainer).(*dig.Container)
	if !exists {
		return errors.New("container is unset in the context")
	}

	if errs, _ := cli.GetMetadataFromContext(ctx, contextKeyDIProvideErrors).([]error); len(errs) > 0 { //nolint:revive,errcheck // unchecked-type-assertion: we know this type for sure
		return fmt.Errorf("provider error: %v", errors.Join(errs...))
	}

	if err := container.Invoke(f, opts...); err != nil {
		return fmt.Errorf("invoker error: %v", dig.RootCause(err))
	}

	return nil
}

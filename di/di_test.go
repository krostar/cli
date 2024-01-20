package clidi

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/krostar/cli"
)

func Test_DI(t *testing.T) {
	type (
		fooA func()
		fooB func()
	)

	t.Run("ok", func(t *testing.T) {
		called := make(map[string]uint)

		ctx := cli.NewContextWithMetadata(context.Background())

		InitializeContainer(ctx)
		AddProvider(ctx, func() fooA { return func() { called["fooA"]++ } })
		AddProvider(ctx, func(a fooA) fooB {
			return func() {
				a()
				called["fooB"]++
			}
		})

		assert.NilError(t, Invoke(ctx, func(b fooB) {
			b()
			called["invoked"]++
		}))
		assert.DeepEqual(t, called, map[string]uint{"fooA": 1, "fooB": 1, "invoked": 1})
	})

	t.Run("calling without metadata", func(t *testing.T) {
		ctx := context.Background()

		InitializeContainer(ctx)
		AddProvider(ctx, func() fooA { return func() {} })
		assert.Error(t, Invoke(ctx, func(fooA) {}), "container is unset in the context")
	})

	t.Run("provider error", func(t *testing.T) {
		ctx := cli.NewContextWithMetadata(context.Background())

		InitializeContainer(ctx)
		AddProvider(ctx, nil)
		assert.ErrorContains(t, Invoke(ctx, func(fooA) {}), "provider error")
	})

	t.Run("invoker error", func(t *testing.T) {
		ctx := cli.NewContextWithMetadata(context.Background())

		InitializeContainer(ctx)
		assert.ErrorContains(t, Invoke(ctx, func(fooA) {}), "invoker error")
	})
}

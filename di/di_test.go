package clidi

import (
	"strings"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"

	"github.com/krostar/cli"
)

func Test_DI(t *testing.T) {
	type (
		fooA func()
		fooB func()
	)

	t.Run("ok", func(t *testing.T) {
		called := make(map[string]uint)

		ctx := cli.NewContextWithMetadata(test.Context(t))

		InitializeContainer(ctx)
		AddProvider(ctx, func() fooA { return func() { called["fooA"]++ } })
		AddProvider(ctx, func(a fooA) fooB {
			return func() {
				a()
				called["fooB"]++
			}
		})

		test.Require(t, Invoke(ctx, func(b fooB) {
			b()
			called["invoked"]++
		}) == nil)
		test.Assert(check.Compare(t, called, map[string]uint{"fooA": 1, "fooB": 1, "invoked": 1}))
	})

	t.Run("calling without metadata", func(t *testing.T) {
		ctx := test.Context(t)

		InitializeContainer(ctx)
		AddProvider(ctx, func() fooA { return func() {} })

		err := Invoke(ctx, func(fooA) {})
		test.Assert(t, err != nil && strings.Contains(err.Error(), "container is unset in the context"))
	})

	t.Run("provider error", func(t *testing.T) {
		ctx := cli.NewContextWithMetadata(test.Context(t))

		InitializeContainer(ctx)
		AddProvider(ctx, nil)

		err := Invoke(ctx, func(fooA) {})
		test.Assert(t, err != nil && strings.Contains(err.Error(), "provider error"))
	})

	t.Run("invoker error", func(t *testing.T) {
		ctx := cli.NewContextWithMetadata(test.Context(t))

		InitializeContainer(ctx)

		err := Invoke(ctx, func(fooA) {})
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invoker error"))
	})
}

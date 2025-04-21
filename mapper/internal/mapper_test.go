package mapper

import (
	"context"
	"errors"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"

	"github.com/krostar/cli"
)

func Test_Context(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		ctx := t.Context()
		test.Assert(t, ctx != Context(new(commandWithAll), ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		ctx := t.Context()
		test.Assert(t, ctx == Context(new(commandSimple), ctx))
	})
}

func Test_ShortDescription(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		test.Assert(t, ShortDescription(new(commandWithAll)) == "short description")
	})

	t.Run("not implemented", func(t *testing.T) {
		test.Assert(t, ShortDescription(new(commandSimple)) == "")
	})
}

func Test_Description(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		test.Assert(t, Description(new(commandWithAll)) == "short description\nlong description")
	})

	t.Run("not implemented", func(t *testing.T) {
		test.Assert(t, Description(new(commandSimple)) == "")
	})
}

func Test_Examples(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		test.Assert(check.Compare(t, Examples(new(commandWithAll)), []string{"example"}))
	})

	t.Run("not implemented", func(t *testing.T) {
		test.Assert(check.Compare(t, Examples(new(commandSimple)), []string(nil)))
	})
}

func Test_Usage(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		test.Assert(t, Usage(new(commandWithAll)) == "usage")
	})

	t.Run("not implemented", func(t *testing.T) {
		test.Assert(t, Usage(new(commandSimple)) == "")
	})
}

func Test_Flags(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		f := Flags(new(commandWithAll))
		test.Require(t, len(f) > 0)
		test.Assert(t, f[0].LongName() == "llong")
	})

	t.Run("not implemented", func(t *testing.T) {
		test.Assert(t, len(Flags(new(commandSimple))) == 0)
	})
}

func Test_PersistentFlags(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		f := PersistentFlags(new(commandWithAll))
		test.Require(t, len(f) > 0)
		test.Assert(t, f[0].LongName() == "plong")
	})

	t.Run("not implemented", func(t *testing.T) {
		test.Assert(t, len(PersistentFlags(new(commandSimple))) == 0)
	})
}

func Test_Hook(t *testing.T) {
	ctx := test.Context(t)

	t.Run("implemented", func(t *testing.T) {
		hook := Hook(new(commandWithAll))
		test.Require(t, hook != nil)
		test.Require(t, hook.BeforeCommandExecution != nil)
		test.Require(t, hook.AfterCommandExecution != nil)

		test.Assert(t, hook.BeforeCommandExecution(ctx) != nil)
		test.Assert(t, hook.AfterCommandExecution(ctx) == nil)
	})

	t.Run("not implemented", func(t *testing.T) {
		hook := Hook(new(commandSimple))
		test.Require(t, hook != nil)
		test.Require(t, hook.BeforeCommandExecution != nil)
		test.Require(t, hook.AfterCommandExecution != nil)

		test.Assert(t, hook.BeforeCommandExecution(ctx) == nil)
		test.Assert(t, hook.AfterCommandExecution(ctx) == nil)
	})
}

func Test_PersistentHook(t *testing.T) {
	ctx := test.Context(t)

	t.Run("implemented", func(t *testing.T) {
		hook := PersistentHook(new(commandWithAll))
		test.Require(t, hook != nil)
		test.Require(t, hook.BeforeFlagsDefinition != nil)
		test.Require(t, hook.BeforeCommandExecution != nil)
		test.Require(t, hook.AfterCommandExecution != nil)

		test.Assert(t, hook.BeforeCommandExecution(ctx) == nil)
		test.Assert(t, hook.BeforeFlagsDefinition(ctx) == nil)
		test.Assert(t, hook.AfterCommandExecution(ctx) == nil)
	})

	t.Run("not implemented", func(t *testing.T) {
		hook := PersistentHook(new(commandSimple))
		test.Require(t, hook != nil)
		test.Require(t, hook.BeforeFlagsDefinition != nil)
		test.Require(t, hook.BeforeCommandExecution != nil)
		test.Require(t, hook.AfterCommandExecution != nil)

		test.Assert(t, hook.BeforeFlagsDefinition(ctx) == nil)
		test.Assert(t, hook.BeforeCommandExecution(ctx) == nil)
		test.Assert(t, hook.AfterCommandExecution(ctx) == nil)
	})
}

type commandSimple struct{}

func (commandSimple) Execute(context.Context, []string, []string) error {
	return nil
}

type commandWithAll struct{}

func (commandWithAll) Execute(context.Context, []string, []string) error {
	return nil
}

func (commandWithAll) Usage() string {
	return "usage"
}

func (commandWithAll) Examples() []string {
	return []string{"example"}
}

func (commandWithAll) Description() string {
	return "short description\nlong description"
}

func (commandWithAll) Context(ctx context.Context) context.Context {
	type key string
	return context.WithValue(ctx, key("foo"), "value")
}

func (commandWithAll) Flags() []cli.Flag {
	var b bool
	return []cli.Flag{cli.NewBuiltinFlag[bool]("llong", "s", &b, "descr")}
}

func (commandWithAll) Hook() *cli.Hook {
	return &cli.Hook{
		BeforeCommandExecution: func(context.Context) error {
			return errors.New("hook")
		},
	}
}

func (commandWithAll) PersistentFlags() []cli.Flag {
	var b bool
	return []cli.Flag{cli.NewBuiltinFlag[bool]("plong", "s", &b, "descr")}
}

func (commandWithAll) PersistentHook() *cli.Hook {
	return &cli.Hook{
		BeforeCommandExecution: func(context.Context) error {
			return errors.New("persistent hook")
		},
	}
}

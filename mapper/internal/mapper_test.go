package mapper

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/krostar/cli"
)

func Test_Context(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		ctx := context.TODO()
		assert.Check(t, ctx != Context(new(commandWithAll), ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		ctx := context.TODO()
		assert.Check(t, ctx == Context(new(commandSimple), ctx))
	})
}

func Test_ShortDescription(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Check(t, ShortDescription(new(commandWithAll)) == "short description")
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Check(t, ShortDescription(new(commandSimple)) == "")
	})
}

func Test_Description(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Check(t, Description(new(commandWithAll)) == "short description\nlong description")
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Check(t, Description(new(commandSimple)) == "")
	})
}

func Test_Examples(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.DeepEqual(t, Examples(new(commandWithAll)), []string{"example"})
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.DeepEqual(t, Examples(new(commandSimple)), []string(nil))
	})
}

func Test_Usage(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Check(t, Usage(new(commandWithAll)) == "usage")
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Check(t, Usage(new(commandSimple)) == "")
	})
}

func Test_Flags(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		f := Flags(new(commandWithAll))
		assert.Check(t, len(f) > 0)
		assert.Check(t, f[0].LongName() == "llong")
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Check(t, len(Flags(new(commandSimple))) == 0)
	})
}

func Test_PersistentFlags(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		f := PersistentFlags(new(commandWithAll))
		assert.Check(t, len(f) > 0)
		assert.Check(t, f[0].LongName() == "plong")
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Check(t, len(PersistentFlags(new(commandSimple))) == 0)
	})
}

func Test_Hook(t *testing.T) {
	ctx := context.Background()

	t.Run("implemented", func(t *testing.T) {
		hook := Hook(new(commandWithAll))
		assert.Assert(t, hook != nil)
		assert.Assert(t, hook.BeforeCommandExecution != nil)
		assert.Assert(t, hook.AfterCommandExecution != nil)

		assert.Error(t, hook.BeforeCommandExecution(ctx), "hook")
		assert.NilError(t, hook.AfterCommandExecution(ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		hook := Hook(new(commandSimple))
		assert.Assert(t, hook != nil)
		assert.Assert(t, hook.BeforeCommandExecution != nil)
		assert.Assert(t, hook.AfterCommandExecution != nil)

		assert.NilError(t, hook.BeforeCommandExecution(ctx))
		assert.NilError(t, hook.AfterCommandExecution(ctx))
	})
}

func Test_PersistentHook(t *testing.T) {
	ctx := context.Background()

	t.Run("implemented", func(t *testing.T) {
		hook := PersistentHook(new(commandWithAll))
		assert.Assert(t, hook != nil)
		assert.Assert(t, hook.BeforeFlagsDefinition != nil)
		assert.Assert(t, hook.BeforeCommandExecution != nil)
		assert.Assert(t, hook.AfterCommandExecution != nil)

		assert.NilError(t, hook.BeforeCommandExecution(ctx))
		assert.NilError(t, hook.BeforeFlagsDefinition(ctx))
		assert.NilError(t, hook.AfterCommandExecution(ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		hook := PersistentHook(new(commandSimple))
		assert.Assert(t, hook != nil)
		assert.Assert(t, hook.BeforeFlagsDefinition != nil)
		assert.Assert(t, hook.BeforeCommandExecution != nil)
		assert.Assert(t, hook.AfterCommandExecution != nil)

		assert.NilError(t, hook.BeforeFlagsDefinition(ctx))
		assert.NilError(t, hook.BeforeCommandExecution(ctx))
		assert.NilError(t, hook.AfterCommandExecution(ctx))
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

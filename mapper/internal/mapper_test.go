package mapper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/cli"
)

func Test_Context(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		ctx := context.TODO()
		assert.NotEqual(t, ctx, Context(new(commandWithAll), ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		ctx := context.TODO()
		assert.Equal(t, ctx, Context(new(commandSimple), ctx))
	})
}

func Test_ShortDescription(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Equal(t, "short description", ShortDescription(new(commandWithAll)))
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Equal(t, "", ShortDescription(new(commandSimple)))
	})
}

func Test_Description(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Equal(t, "short description\nlong description", Description(new(commandWithAll)))
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Equal(t, "", Description(new(commandSimple)))
	})
}

func Test_Examples(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Equal(t, []string{"example"}, Examples(new(commandWithAll)))
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Equal(t, []string(nil), Examples(new(commandSimple)))
	})
}

func Test_Usage(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		assert.Equal(t, "usage", Usage(new(commandWithAll)))
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Equal(t, "", Usage(new(commandSimple)))
	})
}

func Test_Flags(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		f := Flags(new(commandWithAll))
		require.NotEmpty(t, f)
		assert.Equal(t, "llong", f[0].LongName())
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Empty(t, Flags(new(commandSimple)))
	})
}

func Test_PersistentFlags(t *testing.T) {
	t.Run("implemented", func(t *testing.T) {
		f := PersistentFlags(new(commandWithAll))
		require.NotEmpty(t, f)
		assert.Equal(t, "plong", f[0].LongName())
	})

	t.Run("not implemented", func(t *testing.T) {
		assert.Empty(t, PersistentFlags(new(commandSimple)))
	})
}

func Test_Hook(t *testing.T) {
	ctx := context.Background()

	t.Run("implemented", func(t *testing.T) {
		hook := Hook(new(commandWithAll))
		require.NotNil(t, hook)
		assert.NotNil(t, hook.BeforeCommandExecution)
		assert.NotNil(t, hook.AfterCommandExecution)

		err := hook.BeforeCommandExecution(ctx)
		assert.Equal(t, "hook", err.Error())
		assert.NoError(t, hook.AfterCommandExecution(ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		hook := Hook(new(commandSimple))
		require.NotNil(t, hook)
		assert.NotNil(t, hook.BeforeCommandExecution)
		assert.NotNil(t, hook.AfterCommandExecution)

		assert.NoError(t, hook.BeforeCommandExecution(ctx))
		assert.NoError(t, hook.AfterCommandExecution(ctx))
	})
}

func Test_PersistentHook(t *testing.T) {
	ctx := context.Background()

	t.Run("implemented", func(t *testing.T) {
		hook := PersistentHook(new(commandWithAll))
		require.NotNil(t, hook)
		assert.NotNil(t, hook.BeforeFlagsDefinition)
		assert.NotNil(t, hook.BeforeCommandExecution)
		assert.NotNil(t, hook.AfterCommandExecution)

		err := hook.BeforeCommandExecution(ctx)
		assert.Equal(t, "persistent hook", err.Error())
		assert.NoError(t, hook.BeforeFlagsDefinition(ctx))
		assert.NoError(t, hook.AfterCommandExecution(ctx))
	})

	t.Run("not implemented", func(t *testing.T) {
		hook := PersistentHook(new(commandSimple))
		require.NotNil(t, hook)
		assert.NotNil(t, hook.BeforeFlagsDefinition)
		assert.NotNil(t, hook.BeforeCommandExecution)
		assert.NotNil(t, hook.AfterCommandExecution)

		assert.NoError(t, hook.BeforeFlagsDefinition(ctx))
		assert.NoError(t, hook.BeforeCommandExecution(ctx))
		assert.NoError(t, hook.AfterCommandExecution(ctx))
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
	return []cli.Flag{cli.NewFlag[bool]("llong", "s", &b, "descr")}
}

func (commandWithAll) Hook() *cli.Hook {
	return &cli.Hook{
		BeforeCommandExecution: func(ctx context.Context) error {
			return errors.New("hook")
		},
	}
}

func (commandWithAll) PersistentFlags() []cli.Flag {
	var b bool
	return []cli.Flag{cli.NewFlag[bool]("plong", "s", &b, "descr")}
}

func (commandWithAll) PersistentHook() *cli.Hook {
	return &cli.Hook{
		BeforeCommandExecution: func(ctx context.Context) error {
			return errors.New("persistent hook")
		},
	}
}

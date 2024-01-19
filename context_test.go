package cli

import (
	"context"
	"syscall"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_ctxCommand(t *testing.T) {
	{ // check context setup
		ctx := NewCommandContext(context.Background())

		value := ctx.Value(ctxKeyCommand)
		assert.Check(t, value != nil)
	}

	{ // check setting and getting values
		{ // unprepared context
			ctx := context.Background()
			SetInitializedFlagsInContext(ctx, []Flag{nil, nil}, []Flag{nil, nil})
			local, persistent := GetInitializedFlagsFromContext(ctx)
			assert.Check(t, local == nil)
			assert.Check(t, persistent == nil)
		}

		{ // prepared context
			ctx := NewCommandContext(context.Background())
			SetInitializedFlagsInContext(ctx, []Flag{nil, nil}, []Flag{nil})
			local, persistent := GetInitializedFlagsFromContext(ctx)
			assert.Check(t, len(local) == 2)
			assert.Check(t, len(persistent) == 1)
		}
	}
}

func Test_ctxMetadata(t *testing.T) {
	{ // check context setup
		ctx := NewContextWithMetadata(context.Background())

		value := ctx.Value(ctxKeyMetadata)
		assert.Check(t, value != nil)
	}

	{ // check setting and getting values
		{ // unprepared context
			ctx := context.Background()
			SetMetadataInContext(ctx, "key", "value")
			assert.Check(t, GetMetadataFromContext(ctx, "key") == nil)
		}

		{ // prepared context
			ctx := NewContextWithMetadata(context.Background())
			SetMetadataInContext(ctx, "key", "value")
			assert.Check(t, GetMetadataFromContext(ctx, "key").(string) == "value")
		}
	}
}

func Test_NewContextCancelableBySignal(t *testing.T) {
	t.Run("calling cancel func cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		assert.NilError(t, ctx.Err())
		cancel()
		<-ctx.Done()
		assert.ErrorIs(t, ctx.Err(), context.Canceled)
	})

	t.Run("sending provided signal cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		assert.NilError(t, ctx.Err())
		assert.NilError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR1))
		<-ctx.Done()
		assert.ErrorIs(t, ctx.Err(), context.Canceled)
	})

	t.Run("sending unknown signal keeps context intact", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		assert.NilError(t, ctx.Err())
		assert.NilError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR2))
		assert.NilError(t, ctx.Err())
	})
}

package cli

import (
	"context"
	"errors"
	"syscall"
	"testing"

	"github.com/krostar/test"
)

func Test_ctxCommand(t *testing.T) {
	{ // check context setup
		ctx := NewCommandContext(test.Context(t))

		value := ctx.Value(ctxKeyCommand)
		test.Assert(t, value != nil)
	}

	{ // check setting and getting values
		{ // unprepared context
			ctx := test.Context(t)
			SetInitializedFlagsInContext(ctx, []Flag{nil, nil}, []Flag{nil, nil})
			local, persistent := GetInitializedFlagsFromContext(ctx)
			test.Assert(t, local == nil)
			test.Assert(t, persistent == nil)
		}

		{ // prepared context
			ctx := NewCommandContext(test.Context(t))
			SetInitializedFlagsInContext(ctx, []Flag{nil, nil}, []Flag{nil})
			local, persistent := GetInitializedFlagsFromContext(ctx)
			test.Assert(t, len(local) == 2)
			test.Assert(t, len(persistent) == 1)
		}
	}
}

func Test_ctxMetadata(t *testing.T) {
	{ // check context setup
		ctx := NewContextWithMetadata(test.Context(t))

		value := ctx.Value(ctxKeyMetadata)
		test.Assert(t, value != nil)
	}

	{ // check setting and getting values
		{ // unprepared context
			ctx := test.Context(t)
			SetMetadataInContext(ctx, "key", "value")
			test.Assert(t, GetMetadataFromContext(ctx, "key") == nil)
		}

		{ // prepared context
			ctx := NewContextWithMetadata(test.Context(t))
			SetMetadataInContext(ctx, "key", "value")
			test.Assert(t, GetMetadataFromContext(ctx, "key").(string) == "value")
		}
	}
}

func Test_NewContextCancelableBySignal(t *testing.T) {
	t.Run("calling cancel func cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		test.Assert(t, ctx.Err() == nil)
		cancel()
		<-ctx.Done()
		test.Assert(t, errors.Is(ctx.Err(), context.Canceled))
	})

	t.Run("sending provided signal cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		test.Assert(t, ctx.Err() == nil)
		test.Assert(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR1) == nil)
		<-ctx.Done()
		test.Assert(t, errors.Is(ctx.Err(), context.Canceled))
	})

	t.Run("sending unknown signal keeps context intact", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		test.Assert(t, ctx.Err() == nil)
		test.Assert(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR2) == nil)
		test.Assert(t, ctx.Err() == nil)
	})
}

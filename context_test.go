package cli

import (
	"context"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ctxMetadata(t *testing.T) {
	{ // check context setup
		ctx := NewContextWithMetadata(context.Background())

		value := ctx.Value(ctxKeyMetadata)
		require.NotNil(t, value)
		require.IsType(t, make(ctxMetadata), value)
	}

	{ // check setting and getting values
		{ // unprepared context
			ctx := context.Background()
			SetMetadataInContext(ctx, "key", "value")
			assert.Nil(t, GetMetadataFromContext(ctx, "key"))
		}

		{ // prepared context
			ctx := NewContextWithMetadata(context.Background())
			SetMetadataInContext(ctx, "key", "value")
			assert.Equal(t, "value", GetMetadataFromContext(ctx, "key").(string))
		}
	}
}

func Test_NewContextCancelableBySignal(t *testing.T) {
	t.Run("calling cancel func cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		require.NoError(t, ctx.Err())
		cancel()
		<-ctx.Done()
		require.Error(t, ctx.Err())
	})

	t.Run("sending provided signal cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		require.NoError(t, ctx.Err())
		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR1))
		<-ctx.Done()
		require.Error(t, ctx.Err())
	})

	t.Run("sending unknown signal keeps context intact", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		require.NoError(t, ctx.Err())
		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR2))
		require.NoError(t, ctx.Err())
	})
}

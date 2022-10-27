package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bufferThatCloses struct {
	bytes.Buffer
	closed bool
}

func (l *bufferThatCloses) Close() error {
	l.closed = true
	return nil
}
func (l *bufferThatCloses) Closed() bool { return l.closed }

func Test_Exit(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		var (
			exitStatus  *int
			exitMessage bufferThatCloses
		)

		Exit(context.Background(), nil,
			WithExitFunc(func(status int) {
				exitStatus = &status
			}),
			WithExitLoggerFunc(func(context.Context) io.WriteCloser {
				return &exitMessage
			}),
		)

		require.NotNil(t, exitStatus)
		assert.Equal(t, 0, *exitStatus)
		assert.Empty(t, exitMessage.String())
		assert.True(t, exitMessage.Closed())
	})

	t.Run("error", func(t *testing.T) {
		t.Run("without additional behavior", func(t *testing.T) {
			var (
				exitStatus  *int
				exitMessage bufferThatCloses
			)

			Exit(context.Background(), errors.New("boom"),
				WithExitFunc(func(status int) {
					exitStatus = &status
				}),
				WithExitLoggerFunc(func(context.Context) io.WriteCloser {
					return &exitMessage
				}),
			)

			require.NotNil(t, exitStatus)
			assert.Equal(t, 255, *exitStatus)
			assert.Equal(t, "boom\n", exitMessage.String())
			assert.True(t, exitMessage.Closed())
		})

		t.Run("with custom status", func(t *testing.T) {
			var (
				exitStatus  *int
				exitMessage bufferThatCloses
			)

			Exit(context.Background(), NewErrorWithExitStatus(errors.New("boom"), 42),
				WithExitFunc(func(status int) {
					exitStatus = &status
				}),
				WithExitLoggerFunc(func(context.Context) io.WriteCloser {
					return &exitMessage
				}),
			)

			require.NotNil(t, exitStatus)
			assert.Equal(t, 42, *exitStatus)
			assert.Equal(t, "boom\n", exitMessage.String())
			assert.True(t, exitMessage.Closed())
		})
	})
}

func Test_ExitOption(t *testing.T) {
	o := new(exitOptions)

	WithExitFunc(os.Exit)(o)
	assert.NotNil(t, o.exitFunc)

	WithExitLoggerFunc(getExitLoggerFromMetadata)(o)
	assert.NotNil(t, o.getLoggerFunc)
}

func Test_loggerInMetadata(t *testing.T) {
	t.Run("get a logger even if none is previously set", func(t *testing.T) {
		assert.NotNil(t, getExitLoggerFromMetadata(context.Background()))
	})

	t.Run("set a logger", func(t *testing.T) {
		ctx := NewContextWithMetadata(context.Background())
		logger := new(bufferThatCloses)
		SetExitLoggerInMetadata(ctx, logger)
		assert.Equal(t, logger, getExitLoggerFromMetadata(ctx))
	})
}

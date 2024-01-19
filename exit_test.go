package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"gotest.tools/v3/assert"
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

		assert.Assert(t, exitStatus != nil)
		assert.Check(t, *exitStatus == 0)
		assert.Check(t, exitMessage.String() == "")
		assert.Check(t, exitMessage.Closed())
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

			assert.Assert(t, exitStatus != nil)
			assert.Check(t, *exitStatus == 255)
			assert.Check(t, exitMessage.String() == "boom\n")
			assert.Check(t, exitMessage.Closed())
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

			assert.Assert(t, exitStatus != nil)
			assert.Check(t, *exitStatus == 42)
			assert.Check(t, exitMessage.String() == "boom\n")
			assert.Check(t, exitMessage.Closed())
		})
	})
}

func Test_ExitOption(t *testing.T) {
	o := new(exitOptions)

	WithExitFunc(os.Exit)(o)
	assert.Assert(t, o.exitFunc != nil)

	WithExitLoggerFunc(getExitLoggerFromMetadata)(o)
	assert.Assert(t, o.getLoggerFunc != nil)
}

func Test_loggerInMetadata(t *testing.T) {
	t.Run("get a logger even if none is previously set", func(t *testing.T) {
		assert.Assert(t, getExitLoggerFromMetadata(context.Background()) != nil)
	})

	t.Run("set a logger", func(t *testing.T) {
		ctx := NewContextWithMetadata(context.Background())
		logger := new(bufferThatCloses)
		SetExitLoggerInMetadata(ctx, logger)
		assert.Check(t, getExitLoggerFromMetadata(ctx) == logger)
	})
}

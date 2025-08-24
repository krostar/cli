package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/krostar/test"
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

		Exit(test.Context(t), nil,
			WithExitFunc(func(status int) {
				exitStatus = &status
			}),
			WithExitLoggerFunc(func(context.Context) io.WriteCloser {
				return &exitMessage
			}),
		)

		test.Require(t, exitStatus != nil)
		test.Assert(t, *exitStatus == 0)
		test.Assert(t, exitMessage.String() == "")
		test.Assert(t, exitMessage.Closed())
	})

	t.Run("error", func(t *testing.T) {
		t.Run("without additional behavior", func(t *testing.T) {
			var (
				exitStatus  *int
				exitMessage bufferThatCloses
			)

			Exit(test.Context(t), errors.New("boom"),
				WithExitFunc(func(status int) {
					exitStatus = &status
				}),
				WithExitLoggerFunc(func(context.Context) io.WriteCloser {
					return &exitMessage
				}),
			)

			test.Require(t, exitStatus != nil)
			test.Assert(t, *exitStatus == 255)
			test.Assert(t, exitMessage.String() == "boom\n")
			test.Assert(t, exitMessage.Closed())
		})

		t.Run("with custom status", func(t *testing.T) {
			var (
				exitStatus  *int
				exitMessage bufferThatCloses
			)

			Exit(test.Context(t), NewErrorWithExitStatus(errors.New("boom"), 42),
				WithExitFunc(func(status int) {
					exitStatus = &status
				}),
				WithExitLoggerFunc(func(context.Context) io.WriteCloser {
					return &exitMessage
				}),
			)

			test.Require(t, exitStatus != nil)
			test.Assert(t, *exitStatus == 42)
			test.Assert(t, exitMessage.String() == "boom\n")
			test.Assert(t, exitMessage.Closed())
		})
	})
}

func Test_ExitOption(t *testing.T) {
	o := new(exitOptions)

	WithExitFunc(os.Exit)(o)
	test.Require(t, o.exitFunc != nil)

	WithExitLoggerFunc(getExitLoggerFromMetadata)(o)
	test.Require(t, o.getLoggerFunc != nil)
}

func Test_loggerInMetadata(t *testing.T) {
	t.Run("get a logger even if none is previously set", func(t *testing.T) {
		test.Require(t, getExitLoggerFromMetadata(test.Context(t)) != nil)
	})

	t.Run("set a logger", func(t *testing.T) {
		ctx := NewContextWithMetadata(test.Context(t))
		logger := new(bufferThatCloses)
		SetExitLoggerInMetadata(ctx, logger)
		test.Assert(t, getExitLoggerFromMetadata(ctx) == logger)
	})
}

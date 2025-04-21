package cli

import (
	"errors"
	"testing"

	"github.com/krostar/test"
)

func Test_ErrorWithHelp(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithHelp(rootErr)

		test.Assert(t, rootErr.Error() == helpErr.Error())
		test.Assert(t, errors.Is(helpErr, rootErr))

		showHelpErr := new(ShowHelpError)
		test.Assert(t, errors.As(helpErr, showHelpErr))
		test.Assert(t, (*showHelpErr).ShowHelp())
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithHelp(nil)
		test.Assert(t, helpErr.Error() == "")

		showHelpErr := new(ShowHelpError)
		test.Require(t, errors.As(helpErr, showHelpErr))
		test.Assert(t, (*showHelpErr).ShowHelp())
	})
}

func Test_ErrorWithExitStatus(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithExitStatus(rootErr, 42)

		test.Assert(t, rootErr.Error() == helpErr.Error())
		test.Assert(t, errors.Is(helpErr, rootErr))

		showHelpErr := new(ExitStatusError)
		test.Require(t, errors.As(helpErr, showHelpErr))
		test.Assert(t, (*showHelpErr).ExitStatus() == uint8(42))
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithExitStatus(nil, 42)
		test.Assert(t, helpErr.Error() == "")

		showHelpErr := new(ExitStatusError)
		test.Require(t, errors.As(helpErr, showHelpErr))
		test.Assert(t, (*showHelpErr).ExitStatus() == uint8(42))
	})
}

package cli

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_ErrorWithHelp(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithHelp(rootErr)

		assert.Check(t, rootErr.Error() == helpErr.Error())
		assert.ErrorIs(t, helpErr, rootErr)

		showHelpErr := new(ShowHelpError)
		assert.Assert(t, errors.As(helpErr, showHelpErr))
		assert.Check(t, (*showHelpErr).ShowHelp())
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithHelp(nil)
		assert.Check(t, helpErr.Error() == "")

		showHelpErr := new(ShowHelpError)
		assert.Assert(t, errors.As(helpErr, showHelpErr))
		assert.Check(t, (*showHelpErr).ShowHelp())
	})
}

func Test_ErrorWithExitStatus(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithExitStatus(rootErr, 42)

		assert.Check(t, rootErr.Error() == helpErr.Error())
		assert.ErrorIs(t, helpErr, rootErr)

		showHelpErr := new(ExitStatusError)
		assert.Assert(t, errors.As(helpErr, showHelpErr))
		assert.Check(t, (*showHelpErr).ExitStatus() == uint8(42))
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithExitStatus(nil, 42)
		assert.Check(t, helpErr.Error() == "")

		showHelpErr := new(ExitStatusError)
		assert.Assert(t, errors.As(helpErr, showHelpErr))
		assert.Check(t, (*showHelpErr).ExitStatus() == uint8(42))
	})
}

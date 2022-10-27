package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorWithHelp(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithHelp(rootErr)

		assert.Equal(t, rootErr.Error(), helpErr.Error())
		assert.True(t, errors.Is(helpErr, rootErr))

		showHelpErr := new(ShowHelpError)
		assert.ErrorAs(t, helpErr, showHelpErr)
		assert.Equal(t, true, (*showHelpErr).ShowHelp())
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithHelp(nil)
		assert.Empty(t, helpErr.Error())

		showHelpErr := new(ShowHelpError)
		assert.ErrorAs(t, helpErr, showHelpErr)
		assert.Equal(t, true, (*showHelpErr).ShowHelp())
	})
}

func Test_ErrorWithExitStatus(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithExitStatus(rootErr, 42)

		assert.Equal(t, rootErr.Error(), helpErr.Error())
		assert.True(t, errors.Is(helpErr, rootErr))

		showHelpErr := new(ExitStatusError)
		assert.ErrorAs(t, helpErr, showHelpErr)
		assert.Equal(t, uint8(42), (*showHelpErr).ExitStatus())
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithExitStatus(nil, 42)
		assert.Empty(t, helpErr.Error())

		showHelpErr := new(ExitStatusError)
		assert.ErrorAs(t, helpErr, showHelpErr)
		assert.Equal(t, uint8(42), (*showHelpErr).ExitStatus())
	})
}

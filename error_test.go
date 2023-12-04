package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ErrorWithHelp(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithHelp(rootErr)

		assert.Equal(t, rootErr.Error(), helpErr.Error())
		require.ErrorIs(t, helpErr, rootErr)

		showHelpErr := new(ShowHelpError)
		require.ErrorAs(t, helpErr, showHelpErr)
		assert.True(t, (*showHelpErr).ShowHelp())
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithHelp(nil)
		assert.Empty(t, helpErr.Error())

		showHelpErr := new(ShowHelpError)
		require.ErrorAs(t, helpErr, showHelpErr)
		assert.True(t, (*showHelpErr).ShowHelp())
	})
}

func Test_ErrorWithExitStatus(t *testing.T) {
	t.Run("non nil error", func(t *testing.T) {
		rootErr := errors.New("boom")
		helpErr := NewErrorWithExitStatus(rootErr, 42)

		assert.Equal(t, rootErr.Error(), helpErr.Error())
		require.ErrorIs(t, helpErr, rootErr)

		showHelpErr := new(ExitStatusError)
		require.ErrorAs(t, helpErr, showHelpErr)
		assert.Equal(t, uint8(42), (*showHelpErr).ExitStatus())
	})

	t.Run("nil error", func(t *testing.T) {
		helpErr := NewErrorWithExitStatus(nil, 42)
		assert.Empty(t, helpErr.Error())

		showHelpErr := new(ExitStatusError)
		require.ErrorAs(t, helpErr, showHelpErr)
		assert.Equal(t, uint8(42), (*showHelpErr).ExitStatus())
	})
}

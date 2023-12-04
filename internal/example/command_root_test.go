package example

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/cli"
)

func Test_CommandRoot_implements_Command(t *testing.T) {
	cmd := new(CommandRoot)
	assert.Implements(t, (*cli.Command)(nil), cmd)
	assert.Implements(t, (*cli.CommandDescription)(nil), cmd)
}

func Test_CommandRoot_Execute(t *testing.T) {
	cmd := new(CommandRoot)

	err := cmd.Execute(context.Background(), nil, nil)
	require.Error(t, err)

	var showUsageErr cli.ShowHelpError
	require.ErrorAs(t, err, &showUsageErr)
	assert.True(t, showUsageErr.ShowHelp())

	var exitStatusErr cli.ExitStatusError
	require.ErrorAs(t, err, &exitStatusErr)
	assert.Equal(t, uint8(0), exitStatusErr.ExitStatus())
}

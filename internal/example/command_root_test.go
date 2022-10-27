package example

import (
	"context"
	"testing"

	"github.com/krostar/cli"
	"github.com/stretchr/testify/assert"
)

func Test_CommandRoot_implements_Command(t *testing.T) {
	cmd := new(CommandRoot)
	assert.Implements(t, (*cli.Command)(nil), cmd)
	assert.Implements(t, (*cli.CommandDescription)(nil), cmd)
}

func Test_CommandRoot_Execute(t *testing.T) {
	cmd := new(CommandRoot)

	err := cmd.Execute(context.Background(), nil, nil)
	assert.Error(t, err)

	var showUsageErr cli.ShowHelpError
	assert.ErrorAs(t, err, &showUsageErr)
	assert.True(t, showUsageErr.ShowHelp())

	var exitStatusErr cli.ExitStatusError
	assert.ErrorAs(t, err, &exitStatusErr)
	assert.Equal(t, uint8(0), exitStatusErr.ExitStatus())
}

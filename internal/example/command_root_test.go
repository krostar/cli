package example

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/krostar/cli"
)

func Test_CommandRoot_Execute(t *testing.T) {
	cmd := new(CommandRoot)

	err := cmd.Execute(context.Background(), nil, nil)
	assert.Check(t, err != nil)

	var showUsageErr cli.ShowHelpError
	assert.Assert(t, errors.As(err, &showUsageErr))
	assert.Check(t, showUsageErr.ShowHelp())

	var exitStatusErr cli.ExitStatusError
	assert.Assert(t, errors.As(err, &exitStatusErr))
	assert.Check(t, exitStatusErr.ExitStatus() == uint8(0))
}

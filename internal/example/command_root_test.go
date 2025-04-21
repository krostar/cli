package example

import (
	"errors"
	"testing"

	"github.com/krostar/test"

	"github.com/krostar/cli"
)

func Test_CommandRoot_Execute(t *testing.T) {
	cmd := new(CommandRoot)

	err := cmd.Execute(test.Context(t), nil, nil)
	test.Assert(t, err != nil)

	var showUsageErr cli.ShowHelpError
	test.Require(t, errors.As(err, &showUsageErr))
	test.Assert(t, showUsageErr.ShowHelp())

	var exitStatusErr cli.ExitStatusError
	test.Require(t, errors.As(err, &exitStatusErr))
	test.Assert(t, exitStatusErr.ExitStatus() == uint8(0))
}

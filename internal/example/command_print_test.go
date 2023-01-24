package example

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krostar/cli"
)

func Test_CommandPrint_implements_Command(t *testing.T) {
	cmd := new(CommandPrint)
	assert.Implements(t, (*cli.Command)(nil), cmd)
	assert.Implements(t, (*cli.CommandUsage)(nil), cmd)
	assert.Implements(t, (*cli.CommandExamples)(nil), cmd)
	assert.Implements(t, (*cli.CommandDescription)(nil), cmd)
}

func Test_CommandPrint_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		output := new(bytes.Buffer)
		cmd := &CommandPrint{Writer: output}

		err := cmd.Execute(ctx, []string{"foo", "bar"}, []string{"foofoo", "barbar"})
		assert.NoError(t, err)
		assert.Equal(t, `args[0] = foo
args[1] = bar
dashedArgs[0] = foofoo
dashedArgs[1] = barbar
`, output.String())
	})

	t.Run("ko", func(t *testing.T) {
		t.Run("bad arguments numbers", func(t *testing.T) {
			for _, test := range []struct {
				args          []string
				errorContains string
			}{
				{
					args:          nil,
					errorContains: "there should be at least 1 arg to print",
				}, {
					args:          []string{},
					errorContains: "there should be at least 1 arg to print",
				}, {
					args:          []string{"a", "b", "c", "d"},
					errorContains: "there should be no more than 3 args to print",
				},
			} {
				cmd := new(CommandPrint)

				err := cmd.Execute(ctx, test.args, nil)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.errorContains)
			}
		})
	})
}

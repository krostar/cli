package example

import (
	"bytes"
	"strings"
	"testing"

	"github.com/krostar/test"
)

func Test_CommandPrint_Execute(t *testing.T) {
	ctx := test.Context(t)

	t.Run("ok", func(t *testing.T) {
		output := new(bytes.Buffer)
		cmd := &CommandPrint{Writer: output}

		test.Require(t, cmd.Execute(ctx, []string{"foo", "bar"}, []string{"foofoo", "barbar"}) == nil)
		test.Assert(t, output.String() == `args[0] = foo
args[1] = bar
dashedArgs[0] = foofoo
dashedArgs[1] = barbar
`)
	})

	t.Run("ko", func(t *testing.T) {
		t.Run("bad arguments numbers", func(t *testing.T) {
			for _, tt := range []struct {
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
				err := cmd.Execute(ctx, tt.args, nil)
				test.Assert(t, err != nil && strings.Contains(err.Error(), tt.errorContains))
			}
		})
	})
}

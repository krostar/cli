package spf13cobra

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/krostar/test"
	"github.com/spf13/cobra"

	"github.com/krostar/cli"
	"github.com/krostar/cli/double"
	"github.com/krostar/cli/mapper"
)

func Test_Execute(t *testing.T) {
	t.Run("applies all options", func(t *testing.T) {
		option1Called, option2Called := false, false
		option1 := func(*cobra.Command) { option1Called = true }
		option2 := func(*cobra.Command) { option2Called = true }

		err := Execute(t.Context(), []string{"root"}, cli.New(double.NewFake()), option1, option2)
		test.Assert(t, err == nil, "%v", err)
		test.Assert(t, option1Called && option2Called)
	})

	t.Run("cli build failed", func(t *testing.T) {
		anError := errors.New("boom")

		c := cli.New(double.NewFake(double.FakeWithPersistentHook(func() *cli.PersistentHook {
			return &cli.PersistentHook{BeforeFlagsDefinition: func(context.Context) error { return anError }}
		})))

		err := Execute(t.Context(), nil, c, ForTest(t))
		test.Require(t, err != nil && errors.Is(err, anError), "%v", err)
		test.Assert(t, strings.Contains(err.Error(), "unable not build cobra command from cli"))
	})

	t.Run("implementation checks", func(t *testing.T) {
		mapper.AssertImplementation(t, func(t *testing.T, args []string, c *cli.CLI) error {
			return Execute(t.Context(), args, c, ForTest(t))
		})
	})
}

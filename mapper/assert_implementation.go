package mapper

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/krostar/test"

	"github.com/krostar/cli"
	"github.com/krostar/cli/double"
)

// AssertImplementation provides a comprehensive test suite for CLI mapper implementations.
// It validates that any mapper correctly handles all aspects of CLI execution including:
// - CLI name resolution and preservation
// - Error handling and propagation (including custom exit statuses and help requests)
// - Flag parsing and inheritance across command hierarchies
// - Hook execution order (persistent and command-specific hooks)
// - Command structure navigation and execution
// - Context propagation through the command chain
//
// The executeFunc parameter should be the mapper's main execution function that takes
// command line arguments and a CLI instance, then executes the appropriate command.
// This allows testing any mapper implementation (e.g., spf13/cobra, urfave/cli, etc.)
// against the same comprehensive test suite.
func AssertImplementation(t *testing.T, executeFunc func(*testing.T, []string, *cli.CLI) error) {
	t.Run("cli name is handled", func(t *testing.T) {
		t.Run("sets CLI name from first argument", func(t *testing.T) {
			c := cli.New(double.NewFake())
			err := executeFunc(t, []string{"cli-name", "arg1"}, c)
			test.Assert(t, err == nil)
			test.Assert(t, c.Name == "cli-name")
		})

		t.Run("keeps existing CLI name", func(t *testing.T) {
			c := cli.New(double.NewFake())
			c.Name = "existing-name"
			err := executeFunc(t, []string{"should-ignore", "arg1"}, c)
			test.Assert(t, err == nil)
			test.Assert(t, c.Name == "existing-name")
		})
	})

	t.Run("execute error behavior", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			spy, spied := double.SpyCLI(cli.New(double.NewFake()))

			err := executeFunc(t, nil, spied)
			test.Assert(t, err == nil, "%v", err)

			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Execute", true)
		})

		anError := errors.New("boom")

		t.Run("regular error", func(t *testing.T) {
			spy, spied := double.SpyCLI(cli.New(double.NewFake(
				double.FakeWithExecute(func(context.Context, []string, []string) error {
					return anError
				}),
			)))

			err := executeFunc(t, nil, spied)
			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Execute", true)

			test.Require(t, err != nil, "%v", err)
			test.Assert(t, errors.Is(err, anError))
		})

		t.Run("custom exit status error", func(t *testing.T) {
			spy, spied := double.SpyCLI(cli.New(double.NewFake(
				double.FakeWithExecute(func(context.Context, []string, []string) error {
					return cli.NewErrorWithExitStatus(anError, 42)
				}),
			)))

			err := executeFunc(t, nil, spied)
			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Execute", true)

			test.Require(t, err != nil, "%v", err)
			test.Assert(t, errors.Is(err, anError))

			var exitErr cli.ExitStatusError
			test.Assert(t, errors.As(err, &exitErr))
			test.Assert(t, exitErr.ExitStatus() == 42, "expected exit status 42, got %d", exitErr.ExitStatus())
		})

		t.Run("custom error with help", func(t *testing.T) {
			spy, spied := double.SpyCLI(cli.New(double.NewFake(
				double.FakeWithExecute(func(context.Context, []string, []string) error {
					return cli.NewErrorWithHelp(anError)
				}),
				double.FakeWithDescription(func() string {
					return "This is a test command description"
				}),
				double.FakeWithExamples(func() []string {
					return []string{
						"example command arg1",
						"example command arg1 arg2",
					}
				}),
				double.FakeWithUsage(func() string {
					return "<arg1> [arg2]"
				}),
			)))

			err := executeFunc(t, nil, spied)
			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Execute", true)

			test.Require(t, err != nil, "%v", err)
			test.Assert(t, errors.Is(err, anError))

			var helpErr cli.ShowHelpError
			test.Assert(t, errors.As(err, &helpErr))
			test.Assert(t, helpErr.ShowHelp())

			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Description", false)
			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Examples", false)
			spy.AssertCommandMethodCalled(t, []string{spied.Name}, "Usage", false)
		})
	})

	t.Run("flags are correctly handled", func(t *testing.T) {
		var (
			flagStr  string
			flagInt  int
			flagBool bool
		)

		err := executeFunc(t, []string{"myapp", "--str", "value", "--int", "42", "--bool"}, cli.
			New(double.NewFake(
				double.FakeWithFlags(func() []cli.Flag {
					return []cli.Flag{
						cli.NewBuiltinFlag("str", "s", &flagStr, "String flag"),
						cli.NewBuiltinFlag("int", "i", &flagInt, "Int flag"),
						cli.NewBuiltinFlag("bool", "b", &flagBool, "Bool flag"),
					}
				}),
			)),
		)
		test.Assert(t, err == nil, "%v", err)
		test.Assert(t, flagStr == "value")
		test.Assert(t, flagInt == 42)
		test.Assert(t, flagBool)
	})

	t.Run("persistent flags are inherited by subcommands", func(t *testing.T) {
		var (
			rootFlag   string
			subFlag    string
			nestedFlag string
		)

		rootCmd := double.NewFake(
			double.FakeWithPersistentFlags(func() []cli.Flag {
				return []cli.Flag{
					cli.NewBuiltinFlag("root-flag", "", &rootFlag, "Root persistent flag"),
				}
			}),
		)

		subCmd := double.NewFake(
			double.FakeWithPersistentFlags(func() []cli.Flag {
				return []cli.Flag{
					cli.NewBuiltinFlag("sub-flag", "", &subFlag, "Sub persistent flag"),
				}
			}),
		)

		nestedCmd := double.NewFake(
			double.FakeWithFlags(func() []cli.Flag {
				return []cli.Flag{
					cli.NewBuiltinFlag("nested-flag", "", &nestedFlag, "Nested flag"),
				}
			}),
		)

		spy, spied := double.SpyCLI(cli.
			New(rootCmd).
			Mount("sub", cli.
				New(subCmd).
				AddCommand("nested", nestedCmd),
			),
		)

		err := executeFunc(t, []string{
			"app", "sub", "nested",
			"--root-flag", "root-value",
			"--sub-flag", "sub-value",
			"--nested-flag", "nested-value",
		}, spied)

		test.Require(t, err == nil, "%v", err)
		spy.AssertCommandMethodCalled(t, []string{spied.Name}, "PersistentFlags", true)
		spy.AssertCommandMethodCalled(t, []string{spied.Name, "sub"}, "PersistentFlags", true)
		spy.AssertCommandMethodCalled(t, []string{spied.Name, "sub", "nested"}, "Flags", true)
		spy.AssertCommandMethodCalled(t, []string{spied.Name, "sub", "nested"}, "Execute", true)

		test.Assert(t, rootFlag == "root-value")
		test.Assert(t, subFlag == "sub-value")
		test.Assert(t, nestedFlag == "nested-value")
	})

	t.Run("hooks execute in correct order", func(t *testing.T) {
		var hookCallOrder []string

		rootCmd := double.NewFake(
			double.FakeWithPersistentHook(func() *cli.PersistentHook {
				return &cli.PersistentHook{
					BeforeFlagsDefinition: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "root:PersistentBeforeFlagsDefinition")
						return nil
					},
					BeforeCommandExecution: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "root:PersistentBeforeCommandExecution")
						return nil
					},
					AfterCommandExecution: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "root:PersistentAfterCommandExecution")
						return nil
					},
				}
			}),
		)

		subCmd := double.NewFake(
			double.FakeWithPersistentHook(func() *cli.PersistentHook {
				return &cli.PersistentHook{
					BeforeFlagsDefinition: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "sub:PersistentBeforeFlagsDefinition")
						return nil
					},
					BeforeCommandExecution: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "sub:PersistentBeforeCommandExecution")
						return nil
					},
					AfterCommandExecution: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "sub:PersistentAfterCommandExecution")
						return nil
					},
				}
			}),
		)

		nestedCmd := double.NewFake(
			double.FakeWithHook(func() *cli.Hook {
				return &cli.Hook{
					BeforeCommandExecution: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "nested:BeforeCommandExecution")
						return nil
					},
					AfterCommandExecution: func(context.Context) error {
						hookCallOrder = append(hookCallOrder, "nested:AfterCommandExecution")
						return nil
					},
				}
			}),
		)

		err := executeFunc(t, []string{"app", "sub", "nested"}, cli.
			New(rootCmd).
			Mount("sub", cli.
				New(subCmd).
				AddCommand("nested", nestedCmd),
			),
		)
		test.Assert(t, err == nil, "%v", err)

		expectedOrder := []string{
			"root:PersistentBeforeFlagsDefinition",
			"sub:PersistentBeforeFlagsDefinition",
			"root:PersistentBeforeCommandExecution",
			"sub:PersistentBeforeCommandExecution",
			"nested:BeforeCommandExecution",
			"nested:AfterCommandExecution",
			"sub:PersistentAfterCommandExecution",
			"root:PersistentAfterCommandExecution",
		}

		test.Assert(t, slices.Equal(hookCallOrder, expectedOrder),
			"Hook execution order incorrect:\nExpected: %v\nActual: %v",
			expectedOrder, hookCallOrder,
		)
	})

	t.Run("command structure and execution", func(t *testing.T) {
		t.Run("correctly executes subcommand", func(t *testing.T) {
			spy, spied := double.SpyCLI(cli.
				New(double.NewFake()).
				Mount("sub1", cli.
					New(double.NewFake()).
					AddCommand("nested1", double.NewFake()),
				).
				Mount("sub2", cli.
					New(double.NewFake()).
					AddCommand("nested2", double.NewFake()),
				),
			)

			err := executeFunc(t, []string{"app", "sub1", "nested1"}, spied)
			test.Require(t, err == nil, "%v", err)

			spy.AssertCommandMethodCalled(t, []string{spied.Name, "sub1", "nested1"}, "Execute", true)
			test.Assert(t, spy.CountCommandMethodCalls([]string{spied.Name}, "Execute") == 0)
			test.Assert(t, spy.CountCommandMethodCalls([]string{spied.Name, "sub1"}, "Execute") == 0)
			test.Assert(t, spy.CountCommandMethodCalls([]string{spied.Name, "sub2"}, "Execute") == 0)
			test.Assert(t, spy.CountCommandMethodCalls([]string{spied.Name, "sub2", "nested2"}, "Execute") == 0)
		})

		t.Run("handles context propagation across command hierarchy", func(t *testing.T) {
			type contextKey string

			const (
				testKey   contextKey = "test-key"
				rootValue contextKey = "root-value"
				subValue  contextKey = "sub-value"
			)

			var (
				capturedRootCtxValue   any
				capturedSubCtxValue    any
				capturedNestedCtxValue any
			)

			rootCmd := double.NewFake(
				double.FakeWithContext(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, testKey, rootValue)
				}),
				double.FakeWithExecute(func(ctx context.Context, _, _ []string) error {
					capturedRootCtxValue = ctx.Value(testKey)
					return nil
				}),
			)

			subCmd := double.NewFake(
				double.FakeWithContext(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, testKey, subValue)
				}),
				double.FakeWithExecute(func(ctx context.Context, _, _ []string) error {
					capturedSubCtxValue = ctx.Value(testKey)
					return nil
				}),
			)

			nestedCmd := double.NewFake(
				double.FakeWithExecute(func(ctx context.Context, _, _ []string) error {
					capturedNestedCtxValue = ctx.Value(testKey)
					return nil
				}),
			)

			createCLI := func() *cli.CLI {
				return cli.
					New(rootCmd).
					Mount("sub", cli.
						New(subCmd).
						AddCommand("nested", nestedCmd),
					)
			}

			err := executeFunc(t, []string{"app"}, createCLI())
			test.Assert(t, err == nil, "%v", err)
			test.Assert(t, capturedRootCtxValue == rootValue)

			err = executeFunc(t, []string{"app", "sub"}, createCLI())
			test.Assert(t, err == nil, "%v", err)
			test.Assert(t, capturedSubCtxValue == subValue)

			err = executeFunc(t, []string{"app", "sub", "nested"}, createCLI())
			test.Assert(t, err == nil, "%v", err)
			test.Assert(t, capturedNestedCtxValue == subValue)
		})
	})
}

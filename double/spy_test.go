package double

import (
	"context"
	"slices"
	"testing"

	"github.com/krostar/test"
	testdouble "github.com/krostar/test/double"

	"github.com/krostar/cli"
)

func Test_SpyCLI(t *testing.T) {
	t.Run("basic spy creation", func(t *testing.T) {
		var called bool

		spy, spiedCLI := SpyCLI(&cli.CLI{
			Name: "test",
			Command: NewFake(FakeWithExecute(func(context.Context, []string, []string) error {
				called = true
				return nil
			})),
		})

		test.Require(t, spy != nil && spiedCLI != nil)
		test.Assert(t, spiedCLI.Name == "test")
		test.Assert(t, spiedCLI.Command != nil && spiedCLI.Command.Execute(t.Context(), nil, nil) == nil)
		test.Assert(t, called)
	})

	t.Run("spy records method calls", func(t *testing.T) {
		spy, spiedCLI := SpyCLI(cli.New(NewFake(
			FakeWithDescription(func() string { return "test description" }),
			FakeWithUsage(func() string { return "test usage" }),
		)))

		cmdDescription, okDescription := spiedCLI.Command.(cli.CommandDescription)
		cmdUsage, okUsage := spiedCLI.Command.(cli.CommandUsage)
		_, okExample := spiedCLI.Command.(cli.CommandExamples)
		test.Require(t, okDescription && okUsage && !okExample)

		test.Assert(t, cmdDescription.Description() == "test description")
		test.Assert(t, cmdUsage.Usage() == "test usage")
		test.Assert(t, spy.CountCommandMethodCalls([]string{spiedCLI.Name}, "Description") == 1)
		test.Assert(t, spy.CountCommandMethodCalls([]string{spiedCLI.Name}, "Usage") == 1)
	})
}

func Test_Spy_ForEachCommandRecords(t *testing.T) {
	spy, spiedCLI := SpyCLI(cli.New(NewFake(
		FakeWithDescription(func() string { return "desc" }),
		FakeWithUsage(func() string { return "usage" }),
	)))

	cmd := spiedCLI.Command
	cmd.(cli.CommandDescription).Description()
	cmd.(cli.CommandUsage).Usage()

	var recordedMethods []string

	spy.ForEachCommandRecords(func(_ []*cli.CLI, record SpyCommandRecord) {
		recordedMethods = append(recordedMethods, record.Method)
	})

	test.Assert(t, slices.Equal(recordedMethods, []string{"Description", "Usage"}))
}

func Test_Spy_CountCommandMethodCalls(t *testing.T) {
	spy, spiedCLI := SpyCLI(cli.New(NewFake(FakeWithDescription(func() string { return "desc" }))))

	cmd := spiedCLI.Command.(cli.CommandDescription)
	test.Assert(t, spy.CountCommandMethodCalls([]string{spiedCLI.Name}, "Description") == 0)
	cmd.Description()
	test.Assert(t, spy.CountCommandMethodCalls([]string{spiedCLI.Name}, "Description") == 1)
	cmd.Description()
	test.Assert(t, spy.CountCommandMethodCalls([]string{spiedCLI.Name}, "Description") == 2)
}

func Test_Spy_DebugMethods(t *testing.T) {
	spy, spiedCLI := SpyCLI(cli.New(NewFake()).AddCommand("sub", NewFake()))
	spiedCLI.Name = "root"

	_ = spiedCLI.Command.Execute(t.Context(), nil, nil)
	_ = spiedCLI.SubCommands[0].Command.Execute(t.Context(), nil, nil)

	spiedT := testdouble.NewSpy(testdouble.NewFake())
	spy.DebugMethods(spiedT)
	spiedT.ExpectTestToPass(t)
	spiedT.ExpectLogsToContain(t, "[root] Execute called", "[root.sub] Execute called")
}

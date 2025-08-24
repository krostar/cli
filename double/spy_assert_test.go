package double

import (
	"errors"
	"testing"

	"github.com/krostar/test"
	testdouble "github.com/krostar/test/double"

	"github.com/krostar/cli"
)

func Test_Spy_AssertCommandMethodCalled(t *testing.T) {
	spy, spiedCLI := SpyCLI(cli.New(NewFake(FakeWithDescription(func() string { return "test" }))))
	spiedCLI.Name = "test"

	t.Run("not called", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandMethodCalled(spiedT, []string{spiedCLI.Name}, "Description", false)
		spiedT.ExpectTestToFail(t)
		spiedT.ExpectLogsToContain(t, "count is less than 1", "expected test method Description to be called at least once but was not called")
	})

	spiedCLI.Command.(cli.CommandDescription).Description()

	t.Run("called once", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandMethodCalled(spiedT, []string{spiedCLI.Name}, "Description", false)
		spiedT.ExpectTestToPass(t)
		spy.AssertCommandMethodCalled(spiedT, []string{spiedCLI.Name}, "Description", true)
		spiedT.ExpectTestToPass(t)
	})

	spiedCLI.Command.(cli.CommandDescription).Description()

	t.Run("at least once", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandMethodCalled(spiedT, []string{spiedCLI.Name}, "Description", false)
		spiedT.ExpectTestToPass(t)

		spy.AssertCommandMethodCalled(spiedT, []string{spiedCLI.Name}, "Description", true)
		spiedT.ExpectTestToFail(t)
		spiedT.ExpectLogsToContain(t, "count is not equal to 1", "expected test method Description to be called once but called 2 times")
	})
}

func Test_Spy_AssertCommandMethodSequence(t *testing.T) {
	spy, spiedCLI := SpyCLI(cli.New(NewFake(
		FakeWithDescription(func() string { return "test" }),
		FakeWithUsage(func() string { return "usage" }),
		FakeWithExamples(func() []string { return []string{"example"} }),
	)))
	spiedCLI.Name = "test"

	cmd := spiedCLI.Command
	cmd.(cli.CommandDescription).Description()
	cmd.(cli.CommandUsage).Usage()
	cmd.(cli.CommandExamples).Examples()
	cmd.(cli.CommandExamples).Examples()

	t.Run("correct complete sequence", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandMethodSequence(spiedT, []string{spiedCLI.Name}, "Description", "Usage", "Examples", "Examples")
		spiedT.ExpectTestToPass(t)
	})

	t.Run("unmatching number of methods", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandMethodSequence(spiedT, []string{spiedCLI.Name}, "Description", "Usage")
		spiedT.ExpectTestToFail(t)
		spiedT.ExpectLogsToContain(t, "record expectations failed for test: expected 2 methods but found 4")
	})

	t.Run("incorrect sequence", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandMethodSequence(spiedT, []string{spiedCLI.Name}, "Description", "Usage", "Examples", "not-examples")
		spiedT.ExpectTestToFail(t)
		spiedT.ExpectLogsToContain(t, "expected command at position 3 for command test to be method not-examples but found Examples")
	})
}

func Test_Spy_AssertCommandRecords(t *testing.T) {
	spy, spiedCLI := SpyCLI(cli.New(NewFake(
		FakeWithDescription(func() string { return "test" }),
		FakeWithUsage(func() string { return "usage" }),
	)))
	spiedCLI.Name = "test"

	cmd := spiedCLI.Command
	cmd.(cli.CommandDescription).Description()
	cmd.(cli.CommandDescription).Description()
	cmd.(cli.CommandUsage).Usage()

	t.Run("check passed", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandRecords(spiedT, []string{spiedCLI.Name}, func(records []SpyCommandRecord) error {
			test.Require(t, len(records) == 3)
			test.Assert(t, records[0].Method == "Description")
			test.Assert(t, records[1].Method == "Description")
			test.Assert(t, records[2].Method == "Usage")

			return nil
		})
		spiedT.ExpectTestToPass(t)
	})

	t.Run("check failed", func(t *testing.T) {
		spiedT := testdouble.NewSpy(testdouble.NewFake())
		spy.AssertCommandRecords(spiedT, []string{spiedCLI.Name}, func([]SpyCommandRecord) error {
			return errors.New("boom")
		})
		spiedT.ExpectTestToFail(t)
		spiedT.ExpectLogsToContain(t, "record expectations failed for test: boom")
	})
}

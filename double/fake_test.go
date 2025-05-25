package double

import (
	"context"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"

	"github.com/krostar/cli"
)

func Test_Fake(t *testing.T) {
	t.Run("default only implements Execute", func(t *testing.T) {
		f := NewFake()

		test.Assert(check.Not(check.Panics(t, func() {
			test.Assert(t, f.Execute(t.Context(), []string{}, []string{}) == nil)
		}, nil)))

		_, okContext := f.(cli.CommandContext)
		_, okDescription := f.(cli.CommandDescription)
		_, okExamples := f.(cli.CommandExamples)
		_, okFlags := f.(cli.CommandFlags)
		_, okHook := f.(cli.CommandHook)
		_, okPersistentFlags := f.(cli.CommandPersistentFlags)
		_, okPersistentHook := f.(cli.CommandPersistentHook)
		_, okUsage := f.(cli.CommandUsage)

		test.Assert(t, !okContext && !okDescription && !okExamples && !okFlags && !okHook && !okPersistentFlags && !okPersistentHook && !okUsage)
	})

	t.Run("FakeWithContext", func(t *testing.T) {
		ctx := t.Context()

		f := NewFake(FakeWithContext(func(context.Context) context.Context { return ctx }))

		fContext, okContext := f.(cli.CommandContext)
		test.Assert(t, okContext && fContext.Context(t.Context()) == ctx)
	})

	t.Run("FakeWithDescription", func(t *testing.T) {
		description := "hello world"

		f := NewFake(FakeWithDescription(func() string { return description }))

		fDescription, okDescription := f.(cli.CommandDescription)
		test.Assert(t, okDescription && fDescription.Description() == description)
	})

	t.Run("FakeWithExamples", func(t *testing.T) {
		examples := []string{"hello", "world"}

		f := NewFake(FakeWithExamples(func() []string { return examples }))

		fExamples, okExamples := f.(cli.CommandExamples)
		test.Assert(t, okExamples)
		test.Assert(check.Compare(t, fExamples.Examples(), examples))
	})

	t.Run("FakeWithFlags", func(t *testing.T) {
		flags := []cli.Flag{cli.NewBuiltinFlag("long", "s", new(int), "description")}

		f := NewFake(FakeWithFlags(func() []cli.Flag { return flags }))

		fFlags, okFlags := f.(cli.CommandFlags)
		test.Assert(t, okFlags && len(fFlags.Flags()) == len(flags))
	})

	t.Run("FakeWithHook", func(t *testing.T) {
		hook := &cli.Hook{BeforeCommandExecution: func(context.Context) error { return nil }}

		f := NewFake(FakeWithHook(func() *cli.Hook { return hook }))

		fHook, okHook := f.(cli.CommandHook)
		test.Assert(t, okHook && fHook.Hook() == hook)
	})

	t.Run("FakeWithPersistentFlags", func(t *testing.T) {
		flags := []cli.Flag{cli.NewBuiltinFlag("long", "s", new(int), "description")}

		f := NewFake(FakeWithPersistentFlags(func() []cli.Flag { return flags }))

		fPersistentFlags, okPersistentFlags := f.(cli.CommandPersistentFlags)
		test.Assert(t, okPersistentFlags && len(fPersistentFlags.PersistentFlags()) == len(flags))
	})

	t.Run("FakeWithPersistentHook", func(t *testing.T) {
		hook := &cli.PersistentHook{BeforeCommandExecution: func(context.Context) error { return nil }}

		f := NewFake(FakeWithPersistentHook(func() *cli.PersistentHook { return hook }))

		fPersistentHook, okPersistentHook := f.(cli.CommandPersistentHook)
		test.Assert(t, okPersistentHook && fPersistentHook.PersistentHook() == hook)
	})

	t.Run("FakeWithUsage", func(t *testing.T) {
		usage := "hello world"

		f := NewFake(FakeWithUsage(func() string { return usage }))

		fUsage, okUsage := f.(cli.CommandUsage)
		test.Assert(t, okUsage && fUsage.Usage() == usage)
	})
}

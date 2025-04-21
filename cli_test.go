package cli

import (
	"context"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func Test_CLI(t *testing.T) {
	cmd0 := new(command0)
	cmd1 := new(command1)
	cmd2 := new(command2)
	cmd3 := new(command3)
	cmd31 := new(command31)
	cmd4 := new(command4)

	cli := New(cmd0).
		AddCommand("cmd1", cmd1).
		AddCommand("cmd2", cmd2).
		Mount("cmd3", New(cmd3).AddCommand("cmd31", cmd31)).
		Mount("cmd4", New(cmd4))

	test.Assert(check.Compare(t, cli, &CLI{
		Command: cmd0,
		SubCommands: []*CLI{
			{
				Name:    "cmd1",
				Command: cmd1,
			},
			{
				Name:    "cmd2",
				Command: cmd2,
			},
			{
				Name:    "cmd3",
				Command: cmd3,
				SubCommands: []*CLI{
					{
						Name:    "cmd31",
						Command: cmd31,
					},
				},
			},
			{
				Name:    "cmd4",
				Command: cmd4,
			},
		},
	}))
}

type command0 struct{}

func (command0) Execute(context.Context, []string, []string) error { return nil }

type command1 struct{}

func (command1) Execute(context.Context, []string, []string) error { return nil }

type command2 struct{}

func (command2) Execute(context.Context, []string, []string) error { return nil }

type command3 struct{}

func (command3) Execute(context.Context, []string, []string) error { return nil }

type command31 struct{}

func (command31) Execute(context.Context, []string, []string) error { return nil }

type command4 struct{}

func (command4) Execute(context.Context, []string, []string) error { return nil }

package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CLI(t *testing.T) {
	cmd0 := new(command0)
	cmd1 := new(command1)
	cmd2 := new(command2)
	cmd3 := new(command3)
	cmd31 := new(command31)

	cli := NewCommand("cmd0", cmd0).
		AddCommand("cmd1", cmd1).
		AddCommand("cmd2", cmd2).
		Add(NewCommand("cmd3", cmd3).
			AddCommand("cmd31", cmd31))
	assert.Equal(t, &CLI{
		Name:    "cmd0",
		Command: cmd0,
		SubCommands: []*CLI{
			{
				Name:        "cmd1",
				Command:     cmd1,
				SubCommands: nil,
			},
			{
				Name:        "cmd2",
				Command:     cmd2,
				SubCommands: nil,
			},
			{
				Name:    "cmd3",
				Command: cmd3,
				SubCommands: []*CLI{
					{
						Name:        "cmd31",
						Command:     cmd31,
						SubCommands: nil,
					},
				},
			},
		},
	}, cli)
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
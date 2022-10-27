package example

import (
	"context"

	"github.com/krostar/cli"
)

type CommandRoot struct{}

func (CommandRoot) Description() string {
	return `This app is built with the abstraction of any cli backend.
The main goal of this app is to demonstrate how easy it is to switch from one cli backend to another.`
}

func (cmd CommandRoot) Execute(_ context.Context, _ []string, _ []string) error {
	return cli.NewErrorWithExitStatus(cli.NewErrorWithHelp(nil), 0)
}

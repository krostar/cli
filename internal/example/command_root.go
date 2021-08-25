package example

import (
	"context"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
)

type CommandRoot struct{}

func (CommandRoot) Description() string {
	return app.Name() + ` is a cli app built with abstraction of any cli backend.
The main goal of this app is to demonstrate how easy it is to switch from one cli backend to another.`
}

func (cmd CommandRoot) Execute(_ context.Context, _ []string, _ []string) error {
	return cli.ErrorWithExitStatus(cli.ErrorShowHelp(nil), 0)
}

package example

import (
	"context"

	"github.com/krostar/cli"
)

// CommandRoot is the root command for the example CLI application.
// It provides a general description of the application.
type CommandRoot struct{}

// Description returns a description of the example CLI application.
func (CommandRoot) Description() string {
	return `This app is built with the abstraction of any cli backend.
The main goal of this app is to demonstrate how easy it is to switch from one cli backend to another.`
}

// Execute is the main entry point for the root command. It simply returns
// an error that signals the help message to be displayed and sets the exit status to 0.
func (cmd CommandRoot) Execute(_ context.Context, _, _ []string) error {
	return cli.NewErrorWithExitStatus(cli.NewErrorWithHelp(nil), 0)
}

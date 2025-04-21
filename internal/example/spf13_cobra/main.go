package main

import (
	"os"
	"syscall"

	"github.com/krostar/cli"
	"github.com/krostar/cli/internal/example"
	spf13cobra "github.com/krostar/cli/mapper/spf13/cobra"
)

func main() {
	// create a context that can be canceled by SIGINT and SIGTERM signals
	ctx, cancel := cli.NewContextCancelableBySignal(syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// create the CLI with root command and subcommands
	cmd := cli.
		New(new(example.CommandRoot)).
		AddCommand("print", &example.CommandPrint{Writer: os.Stdout})

	// Execute the CLI with spf13/cobra as the backend
	err := spf13cobra.Execute(ctx, os.Args, cmd)

	// Handle exit status and error messages
	cli.Exit(ctx, err)
}

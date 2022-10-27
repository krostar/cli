package main

import (
	"os"
	"syscall"

	"github.com/krostar/cli"
	"github.com/krostar/cli/internal/example"
	spf13cobra "github.com/krostar/cli/mapper/spf13/cobra"
)

func main() {
	ctx, cancel := cli.NewContextCancelableBySignal(syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	cli.Exit(ctx, spf13cobra.Execute(ctx, os.Args, cli.
		NewCommand("my-app", new(example.CommandRoot)).
		AddCommand("print", &example.CommandPrint{Writer: os.Stdout}),
	))
}

package main

import (
	"os"
	"time"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
	"github.com/krostar/cli/internal/example"
	spf13cobra "github.com/krostar/cli/mapper/spf13/cobra"
)

func main() {
	app.Init("app-cli", "1.2.3", time.Now().Format(time.RFC3339))

	ctx, cancel := cli.NewContextCancelableBySignal()
	defer cancel()

	err := spf13cobra.Execute(ctx,
		cli.
			NewCommand("root", &example.CommandRoot{}).
			AddCommand("print", &example.CommandPrint{}),
		os.Args)

	cli.PrintErrorIfAnyAndExit(err)
}

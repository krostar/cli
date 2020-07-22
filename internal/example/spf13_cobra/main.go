package main

import (
	"os"
	"time"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
	"github.com/krostar/cli/internal/example"
	"github.com/krostar/cli/mapper/spf13/cobra"
)

func main() {
	app.Init("app-cli", "1.2.3", time.Now().Format(time.RFC3339))

	ctx, cancel := cli.NewContextCancelableBySignal()
	defer cancel()

	cli.Exit(ctx, cobra.Execute(ctx,
		cli.
			NewCommand("root", &example.CommandRoot{}).
			AddCommand("print", &example.CommandPrint{}),
		os.Args),
	)
}

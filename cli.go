package cli

// CLI stores the settings associated to the CLI.
type CLI struct {
	Name        string
	Command     Command
	SubCommands []*CLI
}

// New creates a new CLI.
func New(cmd Command) *CLI {
	return &CLI{Command: cmd}
}

// AddCommand adds a new subcommand to the CLI.
func (cli *CLI) AddCommand(name string, cmd Command) *CLI {
	cli.SubCommands = append(cli.SubCommands, &CLI{Name: name, Command: cmd})
	return cli
}

// Mount adds a whole new CLI as a subcommand of the CLI.
// Provided name command of the cli is used as sub commands name.
func (cli *CLI) Mount(name string, sub *CLI) *CLI {
	mount := *sub
	mount.Name = name
	cli.SubCommands = append(cli.SubCommands, &mount)
	return cli
}

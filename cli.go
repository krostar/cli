package cli

// CLI stores the settings (name, commands, ...) associated to the CLI.
type CLI struct {
	Name        string
	Command     Command
	SubCommands []*CLI
}

// NewCommand creates a new CLI builder.
func NewCommand(name string, command Command) *CLI {
	return &CLI{
		Name:    name,
		Command: command,
	}
}

// AddCommand adds a new subcommand to the CLI.
func (cli *CLI) AddCommand(name string, cmd Command) *CLI {
	cli.SubCommands = append(cli.SubCommands, NewCommand(name, cmd))
	return cli
}

// Add adds a whole new CLI as a subcommand of the CLI.
func (cli *CLI) Add(sub *CLI) *CLI {
	cli.SubCommands = append(cli.SubCommands, sub)
	return cli
}

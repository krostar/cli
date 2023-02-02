package cli

// CLI stores the settings associated to the CLI.
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
// Name of the root command of the sub cli is used as mount name.
func (cli *CLI) Add(sub *CLI) *CLI {
	cli.SubCommands = append(cli.SubCommands, sub)
	return cli
}

// Mount adds a whole new CLI as a subcommand of the CLI.
// Provided name command of the cli is used as mount name.
func (cli *CLI) Mount(name string, sub *CLI) *CLI {
	ssub := *sub
	ssub.Name = name
	return cli.Add(&ssub)
}

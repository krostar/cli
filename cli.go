// Package cli provides a framework-agnostic way to define command-line interfaces in Go.
//
// The core philosophy of this package is to decouple your application's command-line
// interface logic from specific CLI frameworks like spf13/cobra or urfave/cli. This
// approach offers several benefits:
//
//  1. Testability: Commands and their logic can be easily tested in isolation,
//     without framework dependencies
//
//  2. Flexibility: You can switch between underlying CLI frameworks without
//     significant code changes
//
//  3. Simplicity: The package focuses on core CLI functionality, avoiding
//     unnecessary complexity
//
//  4. Extensibility: Support for custom hooks, flags, and configuration sources
//
// At its core, the package uses a simple Command interface that all CLI commands implement:
//
//	type Command interface {
//		Execute(ctx context.Context, args, dashedArgs []string) error
//	}
//
// Additional interfaces like CommandFlags, CommandDescription, CommandHook, etc.,
// can be implemented to add functionality as needed.
//
// The package also includes a robust configuration system through the cfg subpackage,
// allowing commands to load configuration from multiple sources (environment variables,
// files, command-line flags) with a clear precedence order.
//
// Basic usage:
//
//	cmd := cli.New(rootCommand{}).
//	    AddCommand("sub", subCommand{})
//
//	err := spf13cobra.Execute(context.Background(), os.Args, cmd)
//	cli.Exit(context.Background(), err)
//
// For detailed examples, see the README and the example package.
package cli

// CLI represents a command-line interface with a root command and optional subcommands.
// It provides the structure for building hierarchical CLI applications.
type CLI struct {
	// Name is the name of the command as it appears on the command line.
	// For the root command, this is typically empty as it represents the application itself.
	// For subcommands, it's the command name that will be used to invoke it.
	Name string

	// Command is the command to execute for this CLI.
	// It must implement at least the Command interface, and may optionally
	// implement additional interfaces like CommandFlags, CommandHook, etc.
	Command Command

	// SubCommands is a list of subcommands for this CLI.
	// These are commands that can be invoked under the parent command.
	SubCommands []*CLI
}

// New creates a new CLI with the given command as its root.
// This is the entry point for building a CLI application.
//
// Example:
//
//	rootCmd := New(myRootCommand{})
func New(cmd Command) *CLI {
	return &CLI{Command: cmd}
}

// AddCommand adds a subcommand to the CLI with the given name.
// The name will be used to invoke the command on the command line.
// Returns the CLI instance for chaining method calls.
//
// Example:
//
//	rootCmd := New(myRootCommand{}).
//	    AddCommand("serve", &serveCommand{}).
//	    AddCommand("version", &versionCommand{})
func (cli *CLI) AddCommand(name string, cmd Command) *CLI {
	cli.SubCommands = append(cli.SubCommands, &CLI{Name: name, Command: cmd})
	return cli
}

// Mount adds a pre-configured CLI hierarchy as a subcommand to the current CLI.
// This allows for composing complex CLI structures from simpler ones.
// Unlike AddCommand which adds a single command, Mount adds an entire
// command tree with its own subcommands.
//
// The name parameter sets the name of the mounted CLI's root command.
// Returns the current CLI instance for chaining method calls.
//
// Example:
//
//	// Create a sub-CLI
//	userCLI := New(userRootCommand{}).
//	    AddCommand("list", &listUsersCommand{}).
//	    AddCommand("create", &createUserCommand{})
//
//	// Mount it to the main CLI
//	mainCLI := New(mainRootCommand{}).
//	    Mount("user", userCLI)
func (cli *CLI) Mount(name string, sub *CLI) *CLI {
	mount := *sub
	mount.Name = name
	cli.SubCommands = append(cli.SubCommands, &mount)

	return cli
}

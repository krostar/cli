package cli

import "context"

type (
	// Command is the fundamental interface for all CLI commands.
	Command interface {
		Execute(ctx context.Context, args, dashedArgs []string) error
	}

	// CommandContext allows commands to customize the context passed to
	// their subcommands. This is useful for propagating configuration,
	// dependencies, or other context-specific data down the command tree.
	CommandContext interface {
		Context(ctx context.Context) context.Context
	}

	// CommandDescription allows a command to provide a human-readable
	// description of its purpose. This is used for generating help text.
	// A short description is created from the first description line.
	CommandDescription interface{ Description() string }

	// CommandExamples allows a command to provide usage examples. These
	// examples are displayed in the help text to guide users on how to
	// use the command.
	CommandExamples interface{ Examples() []string }

	// CommandUsage allows a command to specify its argument usage pattern.
	// This helps define how arguments should be passed to the command.
	CommandUsage interface{ Usage() string }

	// CommandFlags allows a command to define command-line flags.
	CommandFlags interface{ Flags() []Flag }

	// CommandPersistentFlags allows a command to define flags that`
	// are inherited by all of its subcommands.
	CommandPersistentFlags interface{ PersistentFlags() []Flag }

	// CommandHook allows a command to define callbacks (hooks) that are
	// executed at specific points in the command's lifecycle. This enables
	// custom behavior before or after command execution.
	CommandHook interface{ Hook() *Hook }

	// CommandPersistentHook allows a command to define persistent hooks
	// that are executed for the command and all of its subcommands.
	CommandPersistentHook interface{ PersistentHook() *PersistentHook }

	// HookFunc defines the signature for hook functions.
	HookFunc func(ctx context.Context) error

	// Hook defines callbacks that are executed during the command lifecycle.
	Hook struct {
		// BeforeCommandExecution is called before the command's Execute method is invoked.
		BeforeCommandExecution HookFunc
		// AfterCommandExecution is called after the command's Execute
		// method has completed (regardless of whether it returned an error).
		AfterCommandExecution HookFunc
	}

	// PersistentHook defines callbacks that are executed for a command and all of its subcommands.
	PersistentHook struct {
		// BeforeFlagsDefinition is called before the command's flags are
		// processed. This is a good place to set up dependencies or
		// perform initialization that affects flag parsing.
		BeforeFlagsDefinition HookFunc
		// BeforeCommandExecution is called before the command's Execute
		// method is invoked (but after flag parsing).
		BeforeCommandExecution HookFunc
		// AfterCommandExecution is called after the command's Execute
		// method has completed (regardless of whether it returned an error).
		AfterCommandExecution HookFunc
	}
)

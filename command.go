package cli

import (
	"context"
)

type (
	// Command defines the minimal interface required to execute a CLI command.
	Command interface {
		Execute(ctx context.Context, args []string, dashedArgs []string) error
	}

	// CommandContext defines a way to propagate a custom context to child commands.
	CommandContext interface {
		Context(context.Context) context.Context
	}

	// CommandDescription defines a way to set a description on the command.
	// A short description is created from the first description line.
	CommandDescription interface{ Description() string }
	// CommandExamples defines a way to set command examples.
	CommandExamples interface{ Examples() []string }
	// CommandUsage defines a way to set the way to use the command.
	CommandUsage interface{ Usage() string }

	// CommandFlags defines the flagValue of the command.
	CommandFlags interface{ Flags() []Flag }
	// CommandPersistentFlags defines the persistent flags of the command.
	CommandPersistentFlags interface{ PersistentFlags() []Flag }

	// CommandHook defines some callback called during command lifecycle.
	CommandHook interface{ Hook() *Hook }
	// CommandPersistentHook defines some persistent callback called during command lifecycle.
	// Differences between CommandHook and CommandPersistentHook is that executed command's
	// hierarchy will also be called.
	CommandPersistentHook interface{ PersistentHook() *Hook }

	// HookFunc defines the hook signature.
	HookFunc func(ctx context.Context) error

	// Hook defines callbacks to add custom behavior to the command lifecycle.
	Hook struct {
		BeforeFlagsDefinition  HookFunc
		BeforeCommandExecution HookFunc
		AfterCommandExecution  HookFunc
	}

	// PersistentHook defines callbacks to add custom behavior to the commands lifecycle.
	PersistentHook struct {
		BeforeCommandExecution HookFunc
		AfterCommandExecution  HookFunc
	}
)

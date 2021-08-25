package mapper

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/krostar/cli"
)

type (
	iDescription     interface{ Description() string }
	iExamples        interface{ Examples() []string }
	iFlags           interface{ Flags() []cli.Flag }
	iHooks           interface{ Hooks() *cli.Hooks }
	iPersistentFlags interface{ PersistentFlags() []cli.Flag }
	iUsage           interface{ Usage() string }
)

// ShortDescription returns the command short description, the first description line.
func ShortDescription(cmd cli.Command) string {
	description := Description(cmd)
	if firstLine, err := bufio.NewReader(strings.NewReader(description)).ReadString('\n'); err == nil || errors.Is(err, io.EOF) {
		return strings.TrimSuffix(firstLine, "\n")
	}
	return description
}

// Description returns the command description.
func Description(cmd cli.Command) string {
	if get, ok := cmd.(iDescription); ok {
		return get.Description()
	}
	return ""
}

// Examples returns the command examples.
func Examples(cmd cli.Command) []string {
	if get, ok := cmd.(iExamples); ok {
		return get.Examples()
	}
	return nil
}

// Flags returns the command flags.
func Flags(cmd cli.Command) []cli.Flag {
	if get, ok := cmd.(iFlags); ok {
		return get.Flags()
	}
	return nil
}

// Hooks returns the command hooks, wth safe default values.
func Hooks(cmd cli.Command) *cli.Hooks {
	var hooks *cli.Hooks

	if get, ok := cmd.(iHooks); ok {
		hooks = get.Hooks()
	}

	if hooks == nil {
		hooks = new(cli.Hooks)
	}

	noopHook := func(context.Context) error { return nil }

	if hooks.BeforeCommandExecution == nil {
		hooks.BeforeCommandExecution = noopHook
	}

	if hooks.AfterCommandExecution == nil {
		hooks.AfterCommandExecution = noopHook
	}

	if hooks.PersistentBeforeCommandExecution == nil {
		hooks.PersistentBeforeCommandExecution = noopHook
	}

	if hooks.PersistentAfterCommandExecution == nil {
		hooks.PersistentAfterCommandExecution = noopHook
	}

	if hooks.BeforeFlagsDefinition == nil {
		hooks.BeforeFlagsDefinition = noopHook
	}

	return hooks
}

// PersistentFlags returns the command persistent flags.
func PersistentFlags(cmd cli.Command) []cli.Flag {
	if get, ok := cmd.(iPersistentFlags); ok {
		return get.PersistentFlags()
	}
	return nil
}

// Usage returns the command usage.
func Usage(cmd cli.Command) string {
	if get, ok := cmd.(iUsage); ok {
		return " " + get.Usage()
	}
	return ""
}

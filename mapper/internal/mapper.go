package mapper

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/krostar/cli"
)

// Context checks whenever command implements interface, with safe default value.
func Context(cmd cli.Command, ctx context.Context) context.Context { //nolint:revive // context-as-argument is voluntarily not respected here to follow other commands signature.
	if get, ok := cmd.(cli.CommandContext); ok {
		return get.Context(ctx)
	}
	return ctx
}

// ShortDescription checks whenever command implements interface, with safe default value.
func ShortDescription(cmd cli.Command) string {
	description := Description(cmd)
	if firstLine, err := bufio.NewReader(strings.NewReader(description)).ReadString('\n'); err == nil || errors.Is(err, io.EOF) {
		return strings.TrimSuffix(firstLine, "\n")
	}
	return description
}

// Description checks whenever command implements interface, with safe default value.
func Description(cmd cli.Command) string {
	if get, ok := cmd.(cli.CommandDescription); ok {
		return get.Description()
	}
	return ""
}

// Examples checks whenever command implements interface, with safe default value.
func Examples(cmd cli.Command) []string {
	if get, ok := cmd.(cli.CommandExamples); ok {
		return get.Examples()
	}
	return nil
}

// Usage checks whenever command implements interface, with safe default value.
func Usage(cmd cli.Command) string {
	if get, ok := cmd.(cli.CommandUsage); ok {
		return get.Usage()
	}
	return ""
}

// Flags checks whenever command implements interface, with safe default value.
func Flags(cmd cli.Command) []cli.Flag {
	if get, ok := cmd.(cli.CommandFlags); ok {
		return get.Flags()
	}
	return nil
}

// PersistentFlags checks whenever command implements interface, with safe default value.
func PersistentFlags(cmd cli.Command) []cli.Flag {
	if get, ok := cmd.(cli.CommandPersistentFlags); ok {
		return get.PersistentFlags()
	}
	return nil
}

// Hook checks whenever command implements interface, with safe default value.
func Hook(cmd cli.Command) *cli.Hook {
	var hooks *cli.Hook

	if get, ok := cmd.(cli.CommandHook); ok {
		hooks = get.Hook()
	}

	if hooks == nil {
		hooks = new(cli.Hook)
	}

	noopHook := func(context.Context) error { return nil }

	if hooks.BeforeCommandExecution == nil {
		hooks.BeforeCommandExecution = noopHook
	}

	if hooks.AfterCommandExecution == nil {
		hooks.AfterCommandExecution = noopHook
	}

	return hooks
}

// PersistentHook checks whenever command implements interface, with safe default value.
func PersistentHook(cmd cli.Command) *cli.PersistentHook {
	var hooks *cli.PersistentHook

	if get, ok := cmd.(cli.CommandPersistentHook); ok {
		hooks = get.PersistentHook()
	}

	if hooks == nil {
		hooks = new(cli.PersistentHook)
	}

	noopHook := func(context.Context) error { return nil }

	if hooks.BeforeCommandExecution == nil {
		hooks.BeforeCommandExecution = noopHook
	}

	if hooks.AfterCommandExecution == nil {
		hooks.AfterCommandExecution = noopHook
	}

	if hooks.BeforeFlagsDefinition == nil {
		hooks.BeforeFlagsDefinition = noopHook
	}

	return hooks
}

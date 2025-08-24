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
// If the provided `cmd` implements the `cli.CommandContext` interface,
// the `Context` method is called to obtain a customized context.
// Otherwise, the original `ctx` is returned.
func Context(cmd cli.Command, ctx context.Context) context.Context { //nolint:revive // context-as-argument is voluntarily not respected here to follow other commands signature.
	if get, ok := cmd.(cli.CommandContext); ok {
		return get.Context(ctx)
	}

	return ctx
}

// ShortDescription checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandDescription` interface,
// its `Description` method is called, and the first line of the description
// is returned as the short description. If an error occurs while reading the
// first line, or if the command does not implement the interface, the full
// description (or an empty string) is returned.
func ShortDescription(cmd cli.Command) string {
	description := Description(cmd)
	if firstLine, err := bufio.NewReader(strings.NewReader(description)).ReadString('\n'); err == nil || errors.Is(err, io.EOF) {
		return strings.TrimSuffix(firstLine, "\n")
	}

	return description
}

// Description checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandDescription` interface,
// its `Description` method is called, and the returned value is used.
// Otherwise, an empty string is returned.
func Description(cmd cli.Command) string {
	if get, ok := cmd.(cli.CommandDescription); ok {
		return get.Description()
	}

	return ""
}

// Examples checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandExamples` interface,
// its `Examples` method is called, and the returned slice of strings is used.
// Otherwise, nil is returned.
func Examples(cmd cli.Command) []string {
	if get, ok := cmd.(cli.CommandExamples); ok {
		return get.Examples()
	}

	return nil
}

// Usage checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandUsage` interface,
// its `Usage` method is called, and the returned value is used.
// Otherwise, an empty string is returned.
func Usage(cmd cli.Command) string {
	if get, ok := cmd.(cli.CommandUsage); ok {
		return get.Usage()
	}

	return ""
}

// Flags checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandFlags` interface,
// its `Flags` method is called, and the returned slice of `cli.Flag` is used.
// Otherwise, nil is returned.
func Flags(cmd cli.Command) []cli.Flag {
	if get, ok := cmd.(cli.CommandFlags); ok {
		return get.Flags()
	}

	return nil
}

// PersistentFlags checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandPersistentFlags` interface,
// its `PersistentFlags` method is called, and the returned slice of `cli.Flag` is used.
// Otherwise, nil is returned.
func PersistentFlags(cmd cli.Command) []cli.Flag {
	if get, ok := cmd.(cli.CommandPersistentFlags); ok {
		return get.PersistentFlags()
	}

	return nil
}

// Hook checks whenever command implements interface, with safe default value.
// If the provided `cmd` implements the `cli.CommandHook` interface,
// its `Hook` method is called, and the returned `*cli.Hook` is used.
// If the command does not implement the interface, or if the returned
// `*cli.Hook` is nil, a new `*cli.Hook` with no-op functions for
// `BeforeCommandExecution` and `AfterCommandExecution` is returned.
// This ensures that the returned `*cli.Hook` always has valid function pointers.
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
// If the provided `cmd` implements the `cli.CommandPersistentHook` interface,
// its `PersistentHook` method is called and the returned value is used.
// If the command doesn't implement the interface or the method returns nil,
// a new `*cli.PersistentHook` is created with no-op functions for
// `BeforeFlagsDefinition`, `BeforeCommandExecution`, and `AfterCommandExecution`.
// This guarantees that the returned `*cli.PersistentHook` is never nil and
// always has valid function pointers.
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

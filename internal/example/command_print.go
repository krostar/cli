package example

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/krostar/cli"
)

type CommandPrint struct {
	Writer io.Writer

	cfg commandPrintConfig
}

type commandPrintConfig struct {
	A bool
	B string
	C []string
}

func (cmd *CommandPrint) Usage() string {
	return "pos args -- dashed args"
}

func (cmd *CommandPrint) Examples() []string {
	return []string{
		"print a",
		"print -- b",
		"print a -- b",
	}
}

func (cmd *CommandPrint) Description() string {
	return "print positional and dashed arguments\n" +
		"print prints at least one and maximum three arguments, and a unlimited number of dashed arguments"
}

func (cmd *CommandPrint) Flags() []cli.Flag {
	return []cli.Flag{
		cli.NewFlag("long-a", "a", &cmd.cfg.A, "displayed when 'a' is a parameter"),
		cli.NewFlag("long-b", "b", &cmd.cfg.B, "displayed when 'b' is a parameter"),
		cli.NewSliceFlag("long-c", "c", &cmd.cfg.C, "displayed when 'c' is a parameter"),
	}
}

func (cmd *CommandPrint) Execute(_ context.Context, args []string, dashedArgs []string) error {
	if len(args) == 0 {
		return cli.NewErrorWithHelp(errors.New("there should be at least 1 arg to print"))
	}
	if len(args) > 3 {
		return errors.New("there should be no more than 3 args to print")
	}

	for i, arg := range args {
		var flag string
		switch arg {
		case "a":
			flag = strconv.FormatBool(cmd.cfg.A)
		case "b":
			flag = cmd.cfg.B
		case "c":
			flag = "[" + strings.Join(cmd.cfg.C, ", ") + "]"
		}

		if flag != "" {
			flag = " (flag=" + flag + ")"
		}

		_, _ = fmt.Fprintf(cmd.Writer, "args[%d] = %s%s\n", i, arg, flag)
	}

	for i, arg := range dashedArgs {
		_, _ = fmt.Fprintf(cmd.Writer, "dashedArgs[%d] = %s\n", i, arg)
	}

	return nil
}

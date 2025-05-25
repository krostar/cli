package spf13cobra

import (
	"testing"

	"github.com/krostar/test"
	"github.com/spf13/pflag"

	"github.com/krostar/cli"
)

func Test_setCobraFlagsFromCLIFlags(t *testing.T) {
	var (
		s string
		i int
		b bool
	)

	flagSet := pflag.NewFlagSet("test", pflag.ExitOnError)
	setCobraFlagsFromCLIFlags(flagSet, []cli.Flag{
		cli.NewBuiltinFlag("string-flag", "s", &s, "Test string flag description"),
		cli.NewBuiltinFlag("int-flag", "", &i, "Test int flag description"),
		cli.NewBuiltinFlag("bool-flag", "b", &b, "Test bool flag description"),
	})

	{ // s
		f := flagSet.Lookup("string-flag")
		test.Assert(t, f != nil, "String flag should be registered")
		test.Assert(t, f.Shorthand == "s", "String flag should have shorthand 's'")
		test.Assert(t, f.Usage == "Test string flag description", "String flag should have correct usage")
		test.Assert(t, f.Value.Set("str") == nil, "String flag should be settable")
		test.Assert(t, s == "str", "String flag should set the variable correctly")
	}

	{ // i
		f := flagSet.Lookup("int-flag")
		test.Assert(t, f != nil, "Int flag should be registered")
		test.Assert(t, f.Shorthand == "", "Int flag should have no shorthand")
		test.Assert(t, f.Usage == "Test int flag description", "Int flag should have correct usage")
		test.Assert(t, f.Value.Set("42") == nil, "Int flag should be settable")
		test.Assert(t, i == 42, "Int flag should set the variable correctly")
	}

	{ // b
		f := flagSet.Lookup("bool-flag")
		test.Assert(t, f != nil, "Bool flag should be registered")
		test.Assert(t, f.Shorthand == "b", "Bool flag should have shorthand 'b'")
		test.Assert(t, f.Usage == "Test bool flag description", "Bool flag should have correct usage")
		test.Assert(t, f.Value.Set("true") == nil, "Bool flag should be settable")
		test.Assert(t, b, "Bool flag should set the variable correctly")
	}
}

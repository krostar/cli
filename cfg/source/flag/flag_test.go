package sourceflag

import (
	"strings"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"

	"github.com/krostar/cli"
)

type configWithFlag struct {
	A string
	B *int
	C map[string]string
	D struct {
		A *string
		B int
	}
	E *struct {
		A string
		B int
	}
}

func Test_Source(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var cfgForFlags configWithFlag
		src := Source[configWithFlag](&cfgForFlags)

		ctx := cli.NewCommandContext(test.Context(t))
		cli.SetInitializedFlagsInContext(ctx,
			[]cli.Flag{
				cli.NewFlag("a", "", cli.NewFlagValuer(&cfgForFlags.C,
					func(s string) (map[string]string, error) {
						m := make(map[string]string)

						lines := strings.Split(s, ",")
						for _, line := range lines {
							split := strings.SplitN(line, ":", 2)
							switch len(split) {
							case 0:
								continue
							case 1:
								m[split[0]] = ""
							case 2:
								m[split[0]] = split[1]
							}
						}

						return m, nil
					},
					func(m map[string]string) string {
						var out []string
						for key, value := range m {
							out = append(out, key+":"+value)
						}
						return strings.Join(out, ",")
					},
				), ""),
			},
			[]cli.Flag{
				cli.NewBuiltinFlag("a", "", &cfgForFlags.A, ""),
				cli.NewBuiltinPointerFlag("E", "", &cfgForFlags.D.A, ""),
			},
		)
		flagsLocal, flagsPersistent := cli.GetInitializedFlagsFromContext(ctx)

		test.Require(t, len(flagsLocal) == 1 && len(flagsPersistent) == 2)
		test.Assert(t, flagsLocal[0].FromString("a:1,b:2") == nil)
		test.Assert(t, flagsPersistent[0].FromString("str") == nil)
		test.Assert(t, flagsPersistent[1].FromString("str") == nil)

		cfg := new(configWithFlag)
		test.Require(t, src(ctx, cfg) == nil)

		test.Assert(check.Compare(t, cfg, &configWithFlag{
			A: "str",
			C: map[string]string{"a": "1", "b": "2"},
			D: struct {
				A *string
				B int
			}{
				A: ptrTo("str"),
			},
		}))
	})

	t.Run("simple structs", func(t *testing.T) {
		type simpleNestedConfig struct {
			Server struct {
				Host string
				Port int
			}
			Database struct {
				DSN string
			}
		}

		var cfgForFlags simpleNestedConfig
		ctx := cli.NewCommandContext(test.Context(t))
		cli.SetInitializedFlagsInContext(ctx,
			[]cli.Flag{
				cli.NewBuiltinFlag("host", "", &cfgForFlags.Server.Host, ""),
				cli.NewBuiltinFlag("port", "", &cfgForFlags.Server.Port, ""),
			},
			[]cli.Flag{
				cli.NewBuiltinFlag("dsn", "", &cfgForFlags.Database.DSN, ""),
			},
		)

		flagsLocal, flagsPersistent := cli.GetInitializedFlagsFromContext(ctx)
		test.Require(t, len(flagsLocal) == 2 && len(flagsPersistent) == 1)

		test.Assert(t, flagsLocal[0].FromString("localhost") == nil)
		test.Assert(t, flagsLocal[1].FromString("8080") == nil)

		test.Assert(t, flagsPersistent[0].FromString("postgres://localhost:5432/db") == nil)

		cfg := new(simpleNestedConfig)
		test.Require(t, Source[simpleNestedConfig](&cfgForFlags)(ctx, cfg) == nil)

		test.Assert(t, cfg.Server.Host == "localhost", "Expected Server.Host to be localhost")
		test.Assert(t, cfg.Server.Port == 8080, "Expected Server.Port to be 8080")
		test.Assert(t, cfg.Database.DSN == "postgres://localhost:5432/db", "Expected Database.DSN to be postgres://localhost:5432/db")
	})

	t.Run("no flags in command", func(t *testing.T) {
		var cfgForFlags configWithFlag
		src := Source[configWithFlag](&cfgForFlags)

		test.Assert(t, src(test.Context(t), new(configWithFlag)) == nil)
	})

	t.Run("no flags set", func(t *testing.T) {
		var cfgForFlags configWithFlag
		src := Source[configWithFlag](&cfgForFlags)

		ctx := cli.NewCommandContext(test.Context(t))
		cli.SetInitializedFlagsInContext(ctx,
			[]cli.Flag{cli.NewBuiltinFlag("a", "", &cfgForFlags.A, "")},
			[]cli.Flag{cli.NewBuiltinFlag("b", "", &cfgForFlags.D.B, "")},
		)

		test.Assert(t, src(ctx, new(configWithFlag)) == nil)
	})

	t.Run("dest is not part of config", func(t *testing.T) {
		var (
			cfgForFlags configWithFlag
			flagDest    string
		)

		src := Source[configWithFlag](&cfgForFlags)
		ctx := cli.NewCommandContext(test.Context(t))
		cli.SetInitializedFlagsInContext(ctx,
			[]cli.Flag{cli.NewBuiltinFlag("a", "", &flagDest, "")},
			[]cli.Flag{cli.NewBuiltinFlag("b", "", &cfgForFlags.D.B, "")},
		)
		flagsLocal, flagsPersistent := cli.GetInitializedFlagsFromContext(ctx)
		test.Require(t, len(flagsLocal) == 1 && len(flagsPersistent) == 1)
		test.Assert(t, flagsLocal[0].FromString("str") == nil)
		test.Assert(t, flagsPersistent[0].FromString("42") == nil)

		err := src(ctx, new(configWithFlag))
		test.Assert(t, err != nil && strings.Contains(err.Error(), "some values where not found, make sure flag values all points to config"))
	})
}

func ptrTo[T any](v T) *T { return &v }

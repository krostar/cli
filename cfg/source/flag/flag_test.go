package sourceflag

import (
	"context"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

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
	t.Run("", func(t *testing.T) {
		var cfgForFlags configWithFlag
		src := Source[configWithFlag](&cfgForFlags)

		ctx := cli.NewCommandContext(context.Background())
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

		assert.NilError(t, flagsLocal[0].FromString("a:1,b:2"))
		assert.NilError(t, flagsPersistent[0].FromString("str"))
		assert.NilError(t, flagsPersistent[1].FromString("str"))

		cfg := new(configWithFlag)
		assert.NilError(t, src(ctx, cfg))

		assert.DeepEqual(t, cfg, &configWithFlag{
			A: "str",
			C: map[string]string{"a": "1", "b": "2"},
			D: struct {
				A *string
				B int
			}{
				A: ptrTo("str"),
			},
		})
	})

	t.Run("no flags in command", func(t *testing.T) {
		var cfgForFlags configWithFlag
		src := Source[configWithFlag](&cfgForFlags)

		assert.NilError(t, src(context.Background(), new(configWithFlag)))
	})

	t.Run("no flags set", func(t *testing.T) {
		var cfgForFlags configWithFlag
		src := Source[configWithFlag](&cfgForFlags)

		ctx := cli.NewCommandContext(context.Background())
		cli.SetInitializedFlagsInContext(ctx,
			[]cli.Flag{cli.NewBuiltinFlag("a", "", &cfgForFlags.A, "")},
			[]cli.Flag{cli.NewBuiltinFlag("b", "", &cfgForFlags.D.B, "")},
		)

		assert.NilError(t, src(ctx, new(configWithFlag)))
	})

	t.Run("dest is not part of config", func(t *testing.T) {
		var (
			cfgForFlags configWithFlag
			flagDest    string
		)

		src := Source[configWithFlag](&cfgForFlags)
		ctx := cli.NewCommandContext(context.Background())
		cli.SetInitializedFlagsInContext(ctx,
			[]cli.Flag{cli.NewBuiltinFlag("a", "", &flagDest, "")},
			[]cli.Flag{cli.NewBuiltinFlag("b", "", &cfgForFlags.D.B, "")},
		)
		flagsLocal, flagsPersistent := cli.GetInitializedFlagsFromContext(ctx)
		assert.NilError(t, flagsLocal[0].FromString("str"))
		assert.NilError(t, flagsPersistent[0].FromString("42"))

		assert.Error(t, src(ctx, new(configWithFlag)), "some values where not find, make sure flag values all points to config")
	})
}

func ptrTo[T any](v T) *T { return &v }

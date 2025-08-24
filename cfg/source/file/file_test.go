package sourcefile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/krostar/test"
	"gopkg.in/yaml.v3"
)

type configWithFile struct {
	A        string `yaml:"awesome"`
	B        string `yaml:"b"`
	Filename string `yaml:"-"`
}

func Test_Source(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		configFile, err := os.CreateTemp(t.TempDir(), "*.yaml")
		test.Require(t, err == nil)
		t.Cleanup(func() { _ = os.Remove(configFile.Name()) })

		_, _ = configFile.WriteString(`
awesome: avalue
`)

		src := Source[configWithFile](
			func(cfg configWithFile) string {
				return cfg.Filename
			},
			func(r io.Reader, cfg *configWithFile) error {
				decoder := yaml.NewDecoder(r)
				decoder.KnownFields(true)

				return decoder.Decode(&cfg)
			},
			false,
		)

		cfg := configWithFile{
			A:        "nota",
			B:        "bvalue",
			Filename: configFile.Name(),
		}
		test.Require(t, src(test.Context(t), &cfg) == nil)
		test.Assert(t, cfg == configWithFile{
			A:        "avalue",
			B:        "bvalue",
			Filename: configFile.Name(),
		})
	})

	t.Run("file not found", func(t *testing.T) {
		configFile, err := os.OpenFile(filepath.Join(t.TempDir(), "perm.yaml"), os.O_CREATE|os.O_EXCL, 0o400)
		test.Require(t, err == nil)
		test.Assert(t, configFile.Close() == nil)
		test.Assert(t, os.Remove(configFile.Name()) == nil)

		t.Run("allow non existing", func(t *testing.T) {
			src := Source[configWithFile](
				func(configWithFile) string { return configFile.Name() },
				func(io.Reader, *configWithFile) error { return errors.New("boom") },
				true,
			)
			test.Assert(t, src(test.Context(t), new(configWithFile)) == nil)
		})

		t.Run("do not allow non existing", func(t *testing.T) {
			src := Source[configWithFile](
				func(configWithFile) string { return configFile.Name() },
				func(io.Reader, *configWithFile) error { return errors.New("boom") },
				false,
			)
			test.Assert(t, errors.Is(src(test.Context(t), new(configWithFile)), os.ErrNotExist))
		})
	})

	t.Run("unable to deserialize", func(t *testing.T) {
		configFile, err := os.CreateTemp(t.TempDir(), "*.yaml")
		test.Require(t, err == nil)
		t.Cleanup(func() { _ = os.Remove(configFile.Name()) })

		expectedErr := errors.New("boom")

		src := Source[configWithFile](
			func(configWithFile) string { return configFile.Name() },
			func(io.Reader, *configWithFile) error { return expectedErr },
			false,
		)

		test.Assert(t, errors.Is(src(test.Context(t), new(configWithFile)), expectedErr))
	})
}

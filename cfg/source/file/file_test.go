package sourcefile

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
	"gotest.tools/v3/assert"
)

type configWithFile struct {
	A        string `yaml:"awesome"`
	B        string `yaml:"b"`
	Filename string `yaml:"-"`
}

func Test_Source(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		configFile, err := os.CreateTemp(t.TempDir(), "*.yaml")
		assert.NilError(t, err)
		t.Cleanup(func() { _ = os.Remove(configFile.Name()) })
		_, _ = io.WriteString(configFile, `
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
		assert.NilError(t, src(context.Background(), &cfg))
		assert.DeepEqual(t, cfg, configWithFile{
			A:        "avalue",
			B:        "bvalue",
			Filename: configFile.Name(),
		})
	})

	t.Run("file not found", func(t *testing.T) {
		configFile, err := os.OpenFile(filepath.Join(t.TempDir(), "perm.yaml"), os.O_CREATE|os.O_EXCL, 0o400)
		assert.NilError(t, err)
		assert.NilError(t, configFile.Close())
		assert.NilError(t, os.Remove(configFile.Name()))

		t.Run("allow non existing", func(t *testing.T) {
			src := Source[configWithFile](
				func(cfg configWithFile) string { return configFile.Name() },
				func(r io.Reader, cfg *configWithFile) error { return errors.New("boom") },
				true,
			)
			assert.NilError(t, src(context.Background(), new(configWithFile)), os.ErrNotExist)
		})

		t.Run("do not allow non existing", func(t *testing.T) {
			src := Source[configWithFile](
				func(cfg configWithFile) string { return configFile.Name() },
				func(r io.Reader, cfg *configWithFile) error { return errors.New("boom") },
				false,
			)
			assert.Check(t, errors.Is(src(context.Background(), new(configWithFile)), os.ErrNotExist))
		})
	})

	t.Run("unable to deserialize", func(t *testing.T) {
		configFile, err := os.CreateTemp(t.TempDir(), "*.yaml")
		assert.NilError(t, err)
		t.Cleanup(func() { _ = os.Remove(configFile.Name()) })

		expectedErr := errors.New("boom")

		src := Source[configWithFile](
			func(cfg configWithFile) string { return configFile.Name() },
			func(r io.Reader, cfg *configWithFile) error { return expectedErr },
			false,
		)

		assert.Check(t, errors.Is(src(context.Background(), new(configWithFile)), expectedErr))
	})
}

package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app, err := New("my-awesome-app", "my-app-version", "2019-12-01T15:44:29Z")
		require.NoError(t, err)

		assert.Equal(t, &App{
			Name:             "my-awesome-app",
			Version:          "my-app-version",
			AlphaNumericName: "myawesomeapp",
			BuiltAt:          time.Date(2019, 12, 01, 15, 44, 29, 0, time.UTC),
		}, app)
	})

	t.Run("wrong build time format", func(t *testing.T) {
		app, err := New("my-awesome-app", "my-app-version", "toto")
		require.Error(t, err)
		assert.Nil(t, app)
	})

	t.Run("empty provided name", func(t *testing.T) {
		app, err := New("", "my-app-version", "2019-12-01T15:44:29Z")
		require.Error(t, err)
		assert.Nil(t, app)
	})

	t.Run("wrong build time format", func(t *testing.T) {
		app, err := New("my-awesome-app", "", "2019-12-01T15:44:29Z")
		require.Error(t, err)
		assert.Nil(t, app)
	})
}

func Test_Copy(t *testing.T) {
	cpy := Copy()
	assert.Equal(t, cpy, app)

	cpy.Name = "not-supposed-to-be-the-same"
	assert.NotEqual(t, cpy.Name, app.Name)
}

func Test_Name(t *testing.T) { assert.NotEmpty(t, Name()) }

func Test_AlphaNumericName(t *testing.T) { assert.NotEmpty(t, AlphaNumericName()) }

func Test_Version(t *testing.T) { assert.NotEmpty(t, Version()) }

func Test_BuiltAt(t *testing.T) { assert.False(t, BuiltAt().IsZero()) }

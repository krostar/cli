package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSliceFlag(t *testing.T) {
	t.Run("no default", func(t *testing.T) {
		var value []int

		flag := NewSliceFlag[int]("longName", "s", &value, "description")
		assert.Equal(t, "longName", flag.LongName())
		assert.Equal(t, "s", flag.ShortName())
		assert.Equal(t, "description", flag.Description())
		assert.Equal(t, "[]int", flag.TypeRepr())

		assert.Error(t, flag.FromString(" 16 ,abc,  18"))

		assert.NoError(t, flag.FromString(" 42 ,  44"))
		repr, err := flag.ToString()
		assert.NoError(t, err)
		assert.Equal(t, "[42,44]", repr)
	})

	t.Run("with defaults", func(t *testing.T) {
		value := []int{42, 44}

		flag := NewSliceFlag[int]("longName", "s", &value, "description")
		assert.Equal(t, "longName", flag.LongName())
		assert.Equal(t, "s", flag.ShortName())
		assert.Equal(t, "description", flag.Description())
		assert.Equal(t, "[]int", flag.TypeRepr())

		assert.Error(t, flag.FromString(" 16 ,abc,  18"))

		assert.NoError(t, flag.FromString(" 46 ,  48"))
		repr, err := flag.ToString()
		assert.NoError(t, err)
		assert.Equal(t, "[42,44,46,48]", repr)
	})
}

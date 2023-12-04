package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewFlag(t *testing.T) {
	var value int

	flag := NewFlag[int]("longName", "s", &value, "description")
	assert.Equal(t, "longName", flag.LongName())
	assert.Equal(t, "s", flag.ShortName())
	assert.Equal(t, "description", flag.Description())
	assert.Equal(t, "int", flag.TypeRepr())

	require.Error(t, flag.FromString("abc"))

	require.NoError(t, flag.FromString("  42 "))
	repr, err := flag.ToString()
	require.NoError(t, err)
	assert.Equal(t, "42", repr)
}

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFlag(t *testing.T) {
	var value int

	flag := NewFlag[int]("longName", "s", &value, "description")
	assert.Equal(t, "longName", flag.LongName())
	assert.Equal(t, "s", flag.ShortName())
	assert.Equal(t, "description", flag.Description())
	assert.Equal(t, "int", flag.TypeRepr())

	assert.Error(t, flag.FromString("abc"))

	assert.NoError(t, flag.FromString("  42 "))
	repr, err := flag.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "42", repr)
}

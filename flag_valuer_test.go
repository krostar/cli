package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCustomFlagValuer(t *testing.T) {
	var d time.Duration

	valuer := NewCustomFlagValuer(&d, time.ParseDuration)

	assert.Equal(t, "time.Duration", valuer.TypeRepr())

	require.Error(t, valuer.FromString("abc"))
	require.NoError(t, valuer.FromString("4s"))
	assert.Equal(t, time.Second*4, d)

	value, err := valuer.ToString()
	require.NoError(t, err)
	assert.Equal(t, "4s", value)
}

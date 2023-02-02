package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewCustomFlagValuer(t *testing.T) {
	var d time.Duration

	valuer := NewCustomFlagValuer(&d, time.ParseDuration)

	assert.Equal(t, "time.Duration", valuer.TypeRepr())

	assert.Error(t, valuer.FromString("abc"))
	assert.NoError(t, valuer.FromString("4s"))
	assert.Equal(t, time.Second*4, d)

	value, err := valuer.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "4s", value)
}

package sequence

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextRange(t *testing.T) {
	sr, err := NextRange("test")
	assert.NoError(t, err)
	assert.NotNil(t, sr)
	assert.Greater(t, sr.max, int64(10))
	assert.Greater(t, sr.min, int64(1))
	assert.Greater(t, sr.value, int64(1))
	assert.False(t, sr.IsOver())
}

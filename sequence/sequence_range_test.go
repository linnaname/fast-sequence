package sequence

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndIncrement(t *testing.T) {
	rs := NewSequenceRange(1, 1000)
	value := rs.GetAndIncrement()
	assert.Equal(t, value, int64(2))
}

func TestIsOver(t *testing.T) {
	rs := NewSequenceRange(1, 3)
	assert.False(t, rs.IsOver())
	rs.GetAndIncrement()
	rs.GetAndIncrement()
	rs.GetAndIncrement()
	assert.True(t, rs.IsOver())
}

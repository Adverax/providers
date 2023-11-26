package interview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache[int, data](NewUnlimitedStorage[int, data](), time.Second)
	go func() {
		e := c.Get(1)
		time.Sleep(1 * time.Second)
		e.Set(&data{1})
	}()

	time.Sleep(100 * time.Millisecond)
	e := c.Get(1)
	val, exists := e.Get()
	require.Equal(t, true, exists)
	assert.Equal(t, 1, val.val)
}

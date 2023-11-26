package interview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTTLStorage_Must(t *testing.T) {
	s := NewTTLStorage[int, data](time.Second)
	s.Set(1, &data{1})
	s.Set(2, &data{2})
	s.Set(3, &data{3})
	s.Set(4, &data{4})
	s.Set(5, &data{5})
	assertTTLStorage(t, s, 1, 2, 3, 4, 5)
	time.Sleep(600 * time.Millisecond)
	s.Set(6, &data{6})
	assertTTLStorage(t, s, 1, 2, 3, 4, 5, 6)
	time.Sleep(800 * time.Millisecond)
	assertTTLStorage(t, s, 6)
}

func assertTTLStorage(t *testing.T, s *TTLStorage[int, data], expected ...int) {
	t.Helper()
	require.Equal(t, len(expected), s.Len())
	for _, v := range expected {
		val, found := s.Get(v)
		require.True(t, found)
		assert.Equal(t, v, val.val)
	}
}

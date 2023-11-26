package interview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type data struct {
	val int
}

func TestHotStorage_Must(t *testing.T) {
	s := NewHotStorage[int, data](5)
	s.Set(1, &data{1})
	s.Set(2, &data{2})
	s.Set(3, &data{3})
	s.Set(4, &data{4})
	s.Set(5, &data{5})
	assertHotStorage(t, s, 1, 2, 3, 4, 5)
	s.Set(6, &data{6})
	assertHotStorage(t, s, 2, 3, 4, 5, 6)
	s.Set(7, &data{7})
	assertHotStorage(t, s, 3, 4, 5, 6, 7)
}

func assertHotStorage(t *testing.T, s *HotStorage[int, data], expected ...int) {
	t.Helper()
	require.Equal(t, len(expected), s.Len())
	for _, v := range expected {
		val, found := s.Get(v)
		require.True(t, found)
		assert.Equal(t, v, val.val)
	}
	item := s.last
	for _, v := range expected {
		require.NotNil(t, item)
		assert.Equal(t, v, item.key)
		item = item.prev
	}
	item = s.first
	for i := len(expected) - 1; i >= 0; i-- {
		require.NotNil(t, item)
		assert.Equal(t, expected[i], item.key)
		item = item.next
	}
}

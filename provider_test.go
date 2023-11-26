package interview

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestProvider(t *testing.T) {
	c := NewCache[int, data](NewUnlimitedStorage[int, data](), time.Second)
	f := &fetcherMock[int, data]{
		data: map[int]*data{
			1: {1},
		},
	}
	p := NewProvider[int, data](f, c)

	v, err := p.Get(1)

	require.NoError(t, err)
	assert.Equal(t, 1, v.val)
}

type fetcherMock[K comparable, V any] struct {
	data map[K]*V
}

func (that *fetcherMock[K, V]) Fetch(key K) (*V, error) {
	v, found := that.data[key]
	if !found {
		return nil, errors.New("not found")
	}
	return v, nil
}

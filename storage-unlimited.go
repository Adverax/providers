package interview

type UnlimitedStorage[K comparable, V any] struct {
	items map[K]*V
}

func NewUnlimitedStorage[K comparable, V any]() *UnlimitedStorage[K, V] {
	return &UnlimitedStorage[K, V]{
		items: make(map[K]*V),
	}
}

func (that *UnlimitedStorage[K, V]) Len() int {
	return len(that.items)
}

func (that *UnlimitedStorage[K, V]) Clear() {
	that.items = make(map[K]*V)
}

func (that *UnlimitedStorage[K, V]) Cleanup() []K {
	return []K{}
}

func (that *UnlimitedStorage[K, V]) Get(key K) (*V, bool) {
	v, found := that.items[key]
	return v, found
}

func (that *UnlimitedStorage[K, V]) Set(key K, val *V) {
	that.items[key] = val
}

package interview

import (
	"container/heap"
	"time"
)

type TTLStorage[K comparable, V any] struct {
	items map[K]*ttlItem[K, V]
	queue expirationHeap[K, V]
	ttl   time.Duration
}

func NewTTLStorage[K comparable, V any](
	ttl time.Duration,
) *TTLStorage[K, V] {
	return &TTLStorage[K, V]{
		items: make(map[K]*ttlItem[K, V]),
		queue: make(expirationHeap[K, V], 0),
		ttl:   ttl,
	}
}

func (that *TTLStorage[K, V]) Len() int {
	var count int
	for _, item := range that.items {
		if !item.isExpired() {
			count++
		}
	}
	return count
}

func (that *TTLStorage[K, V]) Get(key K) (*V, bool) {
	item, found := that.items[key]
	if found {
		if !item.isExpired() {
			return item.value, true
		}
		that.truncateTo(item)
	}
	return nil, false
}

func (that *TTLStorage[K, V]) Set(key K, val *V) {
	item, found := that.items[key]
	if found {
		item.value = val
		item.expiration = time.Now().Add(that.ttl)
		return
	}

	item = &ttlItem[K, V]{
		key:        key,
		value:      val,
		expiration: time.Now().Add(that.ttl),
	}
	that.items[key] = item
	heap.Push(&that.queue, item)
}

func (that *TTLStorage[K, V]) Cleanup() []K {
	result := make([]K, 0, 1000)
	for that.queue.Len() != 0 {
		item := heap.Pop(&that.queue).(*ttlItem[K, V])
		if !item.isExpired() {
			heap.Push(&that.queue, item)
			break
		}

		delete(that.items, item.key)
		result = append(result, item.key)
	}
	return result
}

func (that *TTLStorage[K, V]) Clear() {
	that.items = make(map[K]*ttlItem[K, V])
	that.queue = make(expirationHeap[K, V], 0)
}

func (that *TTLStorage[K, V]) truncateTo(item *ttlItem[K, V]) {
	for i, heapItem := range that.queue {
		heap.Remove(&that.queue, i)
		delete(that.items, heapItem.key)
		if heapItem.key == item.key {
			break
		}
	}
}

type expirationHeap[K comparable, V any] []*ttlItem[K, V]

func (h expirationHeap[K, V]) Len() int           { return len(h) }
func (h expirationHeap[K, V]) Less(i, j int) bool { return h[i].expiration.Before(h[j].expiration) }
func (h expirationHeap[K, V]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *expirationHeap[K, V]) Push(x interface{}) {
	*h = append(*h, x.(*ttlItem[K, V]))
}

func (h *expirationHeap[K, V]) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type ttlItem[K comparable, V any] struct {
	key        K
	value      *V
	expiration time.Time
}

func (that *ttlItem[K, V]) isExpired() bool {
	return time.Now().After(that.expiration)
}

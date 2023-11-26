package interview

type hotItem[K comparable, V any] struct {
	key   K
	value *V
	prev  *hotItem[K, V]
	next  *hotItem[K, V]
}

type HotStorage[K comparable, V any] struct {
	items    map[K]*hotItem[K, V]
	first    *hotItem[K, V]
	last     *hotItem[K, V]
	capacity int
	count    int
}

func NewHotStorage[K comparable, V any](
	capacity int,
) *HotStorage[K, V] {
	return &HotStorage[K, V]{
		items:    make(map[K]*hotItem[K, V]),
		capacity: capacity,
	}
}

func (that *HotStorage[K, V]) Len() int {
	return that.count
}

func (that *HotStorage[K, V]) Clear() {
	that.items = make(map[K]*hotItem[K, V])
	that.first = nil
	that.last = nil
	that.count = 0
}

func (that *HotStorage[K, V]) Cleanup() []K {
	return []K{}
}

func (that *HotStorage[K, V]) Get(key K) (*V, bool) {
	item, found := that.items[key]
	if found {
		that.access(item)
		return item.value, true
	}
	return nil, false
}

func (that *HotStorage[K, V]) Set(key K, val *V) {
	item, found := that.items[key]
	if found {
		item.value = val
		that.access(item)
		return
	}

	that.append(
		&hotItem[K, V]{
			key:   key,
			value: val,
		},
	)
}

func (that *HotStorage[K, V]) access(item *hotItem[K, V]) {
	that.detach(item)
	that.attach(item)
}

func (that *HotStorage[K, V]) detach(item *hotItem[K, V]) {
	if item.prev == nil {
		that.first = item.next
	} else {
		item.prev.next = item.next
	}

	if item.next == nil {
		that.last = item.prev
	} else {
		item.next.prev = item.prev
	}

	that.count--
}

func (that *HotStorage[K, V]) attach(item *hotItem[K, V]) {
	that.count++

	if that.first == nil {
		that.first = item
		that.last = item
		return
	}

	if that.last == that.first {
		that.last.prev = item
	}
	that.first.prev = item
	item.next = that.first
	that.first = item
	if that.count > that.capacity {
		that.remove(that.last.key)
	}
}

func (that *HotStorage[K, V]) append(item *hotItem[K, V]) {
	that.items[item.key] = item
	that.attach(item)
}

func (that *HotStorage[K, V]) remove(key K) {
	item, found := that.items[key]
	if found {
		that.detach(item)
		delete(that.items, key)
	}
}

package interview

import (
	"sync"
	"time"
)

type StorageEntry[V any] interface {
	Get() (*V, bool)
	Set(val *V)
}

type Storage[K comparable, V any] interface {
	Set(key K, val *V)
	Get(key K) (*V, bool)
	Clear()
	Cleanup() []K
	Len() int
}

type Cache[K comparable, V any] struct {
	sync.Mutex
	inFlight map[K]chan struct{}
	storage  Storage[K, V]
	gcPeriod time.Duration
}

func NewCache[K comparable, V any](
	storage Storage[K, V],
	gcPeriod time.Duration,
) *Cache[K, V] {
	c := &Cache[K, V]{
		storage:  storage,
		gcPeriod: gcPeriod,
		inFlight: make(map[K]chan struct{}),
	}
	go c.cleanup()
	return c
}

func (that *Cache[K, V]) Set(key K, val *V) {
	that.Lock()
	defer that.Unlock()

	if val != nil {
		that.storage.Set(key, val)
	}

	that.release(key)
}

func (that *Cache[K, V]) Get(key K) StorageEntry[V] {
	that.Lock()
	defer that.Unlock()

	val, found := that.storage.Get(key)
	if found {
		return &cacheEntry[K, V]{
			cache:    that,
			key:      key,
			val:      val,
			isExists: true,
		}
	}

	ch, inFlight := that.inFlight[key]
	if !inFlight {
		ch = make(chan struct{})
		that.inFlight[key] = ch
	}

	return &cacheEntry[K, V]{
		cache:    that,
		key:      key,
		ch:       ch,
		inFlight: inFlight,
	}
}

func (that *Cache[K, V]) cleanup() {
	for {
		time.Sleep(that.gcPeriod)
		that.cleanup2()
	}
}

func (that *Cache[K, V]) cleanup2() {
	that.Lock()
	defer that.Unlock()

	keys := that.storage.Cleanup()
	for _, key := range keys {
		that.release(key)
	}
}

func (that *Cache[K, V]) release(key K) {
	if ch, inFlight := that.inFlight[key]; inFlight {
		delete(that.inFlight, key)
		close(ch)
	}
}

func (that *Cache[K, V]) Clear() {
	that.Lock()
	defer that.Unlock()

	that.storage.Clear()
	that.inFlight = make(map[K]chan struct{})
}

type cacheEntry[K comparable, V any] struct {
	cache    *Cache[K, V]
	key      K
	val      *V
	ch       chan struct{}
	inFlight bool
	isExists bool
}

func (that *cacheEntry[K, V]) Get() (*V, bool) {
	if that.isExists {
		return that.val, true
	}
	if that.inFlight {
		<-that.ch
		return that.cache.Get(that.key).Get()
	}
	return nil, false
}

func (that *cacheEntry[K, V]) Set(val *V) {
	that.cache.Set(that.key, val)
}

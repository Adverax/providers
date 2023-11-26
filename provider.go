package interview

type Fetcher[K comparable, V any] interface {
	Fetch(key K) (*V, error)
}

type ProviderStorage[K comparable, V any] interface {
	Get(key K) StorageEntry[V]
}

type Provider[K comparable, V any] struct {
	fetcher Fetcher[K, V]
	storage ProviderStorage[K, V]
}

func (that *Provider[K, V]) Get(key K) (*V, error) {
	entry := that.storage.Get(key)
	if val, exists := entry.Get(); exists {
		return val, nil
	}

	val, err := that.fetcher.Fetch(key)

	entry.Set(val)
	return val, err
}

func NewProvider[K comparable, V any](
	fetcher Fetcher[K, V],
	storage ProviderStorage[K, V],
) *Provider[K, V] {
	return &Provider[K, V]{
		fetcher: fetcher,
		storage: storage,
	}
}

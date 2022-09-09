package threadsafe

import (
	"fmt"
	"sync"
)

// ref. MapOf
// https://github.com/SaveTheRbtz/generic-sync-map-go/blob/main/map.go

type ThreadsafeMapSlice[K comparable, V any] struct {
	sync.RWMutex
	items map[K][]V
}

func (tsms *ThreadsafeMapSlice[K, V]) Init(key K) {
	if tsms.items == nil {
		tsms.items = make(map[K][]V)
	}
	if _, ok := tsms.items[key]; !ok {
		tsms.items[key] = make([]V, 0)
	}
}

func (tsms *ThreadsafeMapSlice[K, V]) Append(key K, value V) {
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items[key] = append(tsms.items[key], value)
}

func (tsms *ThreadsafeMapSlice[K, V]) Get(key K) []V {
	tsms.RLock()
	defer tsms.RUnlock()
	return tsms.items[key]
}

func (tsms *ThreadsafeMapSlice[K, V]) Set(key K, value []V) {
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items[key] = value
}

func (tsms *ThreadsafeMapSlice[K, V]) Delete(key K) {
	tsms.Lock()
	defer tsms.Unlock()
	delete(tsms.items, key)
}

func (tsms *ThreadsafeMapSlice[K, V]) DeleteAll() {
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items = make(map[K][]V)
}

func (tsms *ThreadsafeMapSlice[K, V]) Len(key K) int {
	tsms.RLock()
	defer tsms.RUnlock()
	return len(tsms.items[key])
}

func (tsms *ThreadsafeMapSlice[K, V]) Contains(key K) bool {
	tsms.RLock()
	defer tsms.RUnlock()
	_, ok := tsms.items[key]
	return ok
}

func (tsms *ThreadsafeMapSlice[K, V]) Keys() []K {
	tsms.RLock()
	defer tsms.RUnlock()
	keys := make([]K, 0, len(tsms.items))
	for k := range tsms.items {
		keys = append(keys, k)
	}
	return keys
}

func (tsms *ThreadsafeMapSlice[K, V]) KeysAsString() []string {
	tsms.RLock()
	defer tsms.RUnlock()
	var keys []string
	for _, k := range tsms.Keys() {
		keys = append(keys, fmt.Sprint(k))
	}
	return keys
}

func (tsms *ThreadsafeMapSlice[K, V]) Filter(key K, f func(V) bool) (result []V) {
	tsms.RLock()
	defer tsms.RUnlock()
	for _, item := range tsms.items[key] {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

func (tsms *ThreadsafeMapSlice[K, V]) FilterLast(key K, f func(V) bool) (result V) {
	tsms.RLock()
	defer tsms.RUnlock()
	for _, item := range tsms.items[key] {
		if f(item) {
			result = item
		}
	}
	return result
}

func (tsms *ThreadsafeMapSlice[K, V]) FilterDelete(key K, f func(V) bool) {
	var result []V
	for _, item := range tsms.items[key] {
		if !f(item) {
			result = append(result, item)
		}
	}
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items[key] = result
}

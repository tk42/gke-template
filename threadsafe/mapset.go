package threadsafe

import (
	"fmt"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

// ref. MapOf
// https://github.com/SaveTheRbtz/generic-sync-map-go/blob/main/map.go

type ThreadsafeMapSet[K comparable, V comparable] struct {
	sync.RWMutex
	items map[K]mapset.Set[V]
}

func NewThreadsafeMapSet[K comparable, V comparable]() *ThreadsafeMapSet[K, V] {
	tsms := new(ThreadsafeMapSet[K, V])
	if tsms.items == nil {
		tsms.items = make(map[K]mapset.Set[V])
	}
	return tsms
}

func (tsms *ThreadsafeMapSet[K, V]) Append(key K, value V) {
	tsms.Lock()
	defer tsms.Unlock()
	if _, ok := tsms.items[key]; !ok {
		tsms.items[key] = mapset.NewSet[V]()
	}
	tsms.items[key].Add(value)
}

func (tsms *ThreadsafeMapSet[K, V]) Get(key K) mapset.Set[V] {
	tsms.RLock()
	defer tsms.RUnlock()
	return tsms.items[key]
}

func (tsms *ThreadsafeMapSet[K, V]) Set(key K, value mapset.Set[V]) {
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items[key] = value
}

func (tsms *ThreadsafeMapSet[K, V]) Delete(key K, value V) {
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items[key].Remove(value)
}

func (tsms *ThreadsafeMapSet[K, V]) DeleteAll() {
	tsms.Lock()
	defer tsms.Unlock()
	tsms.items = make(map[K]mapset.Set[V])
}

func (tsms *ThreadsafeMapSet[K, V]) Len(key K) int {
	tsms.RLock()
	defer tsms.RUnlock()
	return tsms.items[key].Cardinality()
}

func (tsms *ThreadsafeMapSet[K, V]) Contains(key K) bool {
	tsms.RLock()
	defer tsms.RUnlock()
	_, ok := tsms.items[key]
	return ok
}

func (tsms *ThreadsafeMapSet[K, V]) Keys() []K {
	tsms.RLock()
	defer tsms.RUnlock()
	keys := make([]K, 0, len(tsms.items))
	for k := range tsms.items {
		keys = append(keys, k)
	}
	return keys
}

func (tsms *ThreadsafeMapSet[K, V]) String() []string {
	var values []string
	for _, k := range tsms.Keys() {
		for _, v := range tsms.Get(k).ToSlice() {
			values = append(values, fmt.Sprintf("%v/%v", k, v))
		}
	}
	return values
}

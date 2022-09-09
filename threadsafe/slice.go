package threadsafe

import "sync"

// Slice type that can be safely shared between goroutines
type ThreadsafeSlice[T any] struct {
	sync.RWMutex
	items []T
}

func (cs *ThreadsafeSlice[T]) Append(item T) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = append(cs.items, item)
}

func (cs *ThreadsafeSlice[T]) Range() []T {
	return cs.items
}

func (cs *ThreadsafeSlice[T]) Length() int {
	return len(cs.items)
}

func (cs *ThreadsafeSlice[T]) Filter(f func(T) bool) (result []T) {
	for _, item := range cs.items {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

func (cs *ThreadsafeSlice[T]) FilterLast(f func(T) bool) (result T) {
	for _, item := range cs.items {
		if f(item) {
			result = item
		}
	}
	return result
}

func (cs *ThreadsafeSlice[T]) FilterDelete(f func(T) bool) {
	var result []T
	for _, item := range cs.items {
		if !f(item) {
			result = append(result, item)
		}
	}
	cs.Lock()
	defer cs.Unlock()
	cs.items = result
}

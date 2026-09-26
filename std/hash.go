package std

import "sync"

type Hash[T any] struct {
	data map[string]T
	lock sync.Mutex
}

func NewHash[T any]() *Hash[T] {
	return &Hash[T]{
		data: make(map[string]T),
	}
}

func (h *Hash[T]) Set(key string, value T) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.data[key] = value
}

func (h *Hash[T]) Get(key string) (T, bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	value, ok := h.data[key]

	return value, ok
}

func (h *Hash[T]) Raw() map[string]T {
	return h.data
}

func (h *Hash[T]) IsEmpty() bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	return len(h.data) == 0
}

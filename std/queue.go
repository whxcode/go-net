package std

import "sync"

type Queue[T any] struct {
	lock sync.Mutex
	data []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		data: make([]T, 0),
	}
}

func (q *Queue[T]) Push(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.data = append(q.data, item)
}

func (q *Queue[T]) Pop() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	item := q.data[0]
	q.data = q.data[1:]
	return item, true
}

func (q *Queue[T]) Raw() []T {
	return q.data
}

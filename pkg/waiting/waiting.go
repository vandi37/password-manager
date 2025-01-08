package waiting

import (
	"cmp"
)

type Waiter[K cmp.Ordered, V any] struct {
	queue map[K]chan V
}

func New[K cmp.Ordered, V any]() *Waiter[K, V] {
	return &Waiter[K, V]{
		queue: make(map[K]chan V),
	}
}

func (w *Waiter[K, V]) Add(key K) chan V {
	ch := make(chan V)

	val, ok := w.queue[key]
	if ok {
		close(val)
	}

	w.queue[key] = ch

	return ch
}

func (w *Waiter[K, V]) Check(key K, val V) bool {
	ch, ok := w.queue[key]
	if !ok {
		return false
	}

	ch <- val

	return true
}

func (w *Waiter[K, V]) Remove(key K) bool {
	ch, ok := w.queue[key]
	if !ok {
		return false
	}

	close(ch)
	delete(w.queue, key)

	return true
}

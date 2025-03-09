package waiting

import (
	"cmp"
)

type Waiter[K cmp.Ordered, V any] struct {
	queue map[K]struct {
		ch     chan V
		cancel Cancel
	}
}

func New[K cmp.Ordered, V any]() *Waiter[K, V] {
	return &Waiter[K, V]{
		queue: make(map[K]struct {
			ch     chan V
			cancel Cancel
		}),
	}
}

func (w *Waiter[K, V]) Add(key K) (chan V, Cancel) {
	w.Remove(key)
	ch := make(chan V)
	cancel := NewCancel()

	w.queue[key] = struct {
		ch     chan V
		cancel Cancel
	}{
		ch:     ch,
		cancel: cancel,
	}

	return ch, cancel
}

func (w *Waiter[K, V]) Check(key K, val V) bool {
	ch, ok := w.queue[key]
	if !ok {
		return false
	}

	ch.ch <- val

	return true
}

func (w *Waiter[K, V]) Remove(key K) bool {
	ch, ok := w.queue[key]
	if !ok {
		return false
	}

	close(ch.ch)
	ch.cancel.Canceled()
	delete(w.queue, key)

	return true
}

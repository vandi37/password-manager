package waiting

import (
	"cmp"
)

type Waiter[K cmp.Ordered, V any] struct {
	queue map[K]struct {
		ch       chan V
		cancel   Cancel
		canceled bool
	}
}

func New[K cmp.Ordered, V any]() *Waiter[K, V] {
	return &Waiter[K, V]{
		queue: make(map[K]struct {
			ch       chan V
			cancel   Cancel
			canceled bool
		}),
	}
}

func (w *Waiter[K, V]) Add(key K) (chan V, Cancel) {
	ch := make(chan V)
	cancel := NewCancel()

	val, ok := w.queue[key]
	if ok {
		close(val.ch)
	}

	w.queue[key] = struct {
		ch       chan V
		cancel   Cancel
		canceled bool
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

func (w *Waiter[K, V]) Cancel(key K) bool {
	ch, ok := w.queue[key]
	if !ok || ch.canceled {
		return false
	}

	ch.canceled = true

	ch.cancel.Cancel()
	return true
}

func (w *Waiter[K, V]) Remove(key K) bool {
	ch, ok := w.queue[key]
	if !ok {
		return false
	}

	delete(w.queue, key)
	close(ch.ch)
	close(ch.cancel.cancel)

	return true
}

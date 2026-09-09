// Package mailbox provides asynchronous FIFO message processing.
package mailbox

import "sync"

// Mailbox processes messages in send order without blocking senders on the
// consumer. The queue is unbounded and safe for concurrent senders.
type Mailbox[T any] struct {
	lock    sync.Mutex
	pending []T
	consume func(T)
}

// New returns a mailbox that calls consume for one message at a time.
func New[T any](consume func(T)) *Mailbox[T] {
	return &Mailbox[T]{consume: consume}
}

// Send adds message to the mailbox and starts a worker when the queue was idle.
func (m *Mailbox[T]) Send(message T) {
	m.lock.Lock()
	m.pending = append(m.pending, message)
	start := len(m.pending) == 1
	m.lock.Unlock()

	if start {
		go m.drain()
	}
}

func (m *Mailbox[T]) drain() {
	for {
		m.lock.Lock()
		message := m.pending[0]
		m.lock.Unlock()

		m.consume(message)

		m.lock.Lock()
		var zero T
		m.pending[0] = zero
		m.pending = m.pending[1:]
		if len(m.pending) == 0 {
			m.pending = nil
			m.lock.Unlock()

			return
		}
		m.lock.Unlock()
	}
}

package queue

import (
	"context"
	"errors"
	"sync"
)

var _ Queue[any] = &ConcurrentBlockingQueue[any]{}

type ConcurrentBlockingQueue[T any] struct {
	lock sync.Mutex
	data []any

	head int
	tail int
	size int
}

func NewConcurrentBlockingQueue[T any](size int) *ConcurrentBlockingQueue[T] {
	return &ConcurrentBlockingQueue[T]{
		data: make([]any, size),
	}
}

func (c *ConcurrentBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	if c.IsFull() {
		return errors.New(`queue is full`)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[c.tail] = t
	//c.tail = (c.tail + 1) % c.size
	c.tail++
	if c.tail == len(c.data) {
		c.tail = 0
	}
	c.size++
	return nil
}

func (c *ConcurrentBlockingQueue[T]) Dequeue(ctx context.Context) (t T, err error) {
	if c.IsEmpty() {
		err = errors.New(`queue is empty`)
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	t = c.data[c.head]
	c.head++
	//c.head = (c.head + 1) % c.size
	if c.head == len(c.data) {
		c.head = 0
	}
	c.size--
	return
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	return c.size == 0
}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	return c.size == len(c.data)
}

func (c *ConcurrentBlockingQueue[T]) Len() int {
	return c.size
}

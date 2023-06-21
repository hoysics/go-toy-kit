package queue

import (
	"context"
	"sync"
)

var _ Queue[any] = &ConcurrentBlockingQueue[any]{}

type ConcurrentBlockingQueue[T any] struct {
	lock sync.Mutex
	data []any

	notEmpty chan struct{}
	notFull  chan struct{}

	head int
	tail int
	size int
}

func NewConcurrentBlockingQueue[T any](size int) *ConcurrentBlockingQueue[T] {
	return &ConcurrentBlockingQueue[T]{
		data:     make([]any, size),
		notEmpty: make(chan struct{}, 1),
		notFull:  make(chan struct{}, 1),
	}
}

func (c *ConcurrentBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	c.lock.Lock()
	for c.IsFull() {
		c.lock.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.notFull:
			c.lock.Lock()
		}
	}
	c.data[c.tail] = t
	c.tail++
	if c.tail == len(c.data) {
		c.tail = 0
	}
	c.size++
	// 如果有人在等待有数据 则应该唤醒这个人 但是不能阻塞我自己
	if c.size == 1 {
		// 只有从满变不满发信号
		c.notEmpty <- struct{}{}
	}
	c.lock.Unlock()
	return nil
}

func (c *ConcurrentBlockingQueue[T]) Dequeue(ctx context.Context) (t T, err error) {
	c.lock.Lock()
	for c.IsEmpty() {
		c.lock.Unlock()
		// 注意这两步的时间差
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-c.notEmpty:
			c.lock.Lock()
		}
	}
	t = c.data[c.head]
	c.head++
	if c.head == len(c.data) {
		c.head = 0
	}
	c.size--
	// 如果有人在等待少个数据 则应该唤醒这个人 但是不能阻塞我自己
	if c.size == len(c.data)-1 {
		//只有从空变不空发信号
		c.notFull <- struct{}{}
	}
	c.lock.Unlock()
	return
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.size == 0
}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.size == len(c.data)
}

func (c *ConcurrentBlockingQueue[T]) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.size
}

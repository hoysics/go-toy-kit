package solution

import (
	"context"
	"sync"
)

var _ Queue[any] = &ConcurrentBlockingQueueV1[any]{}

type ConcurrentBlockingQueueV1[T any] struct {
	lock *sync.Mutex
	data []any

	//notEmpty chan struct{}
	//notFull  chan struct{}
	notEmptyCond *cond
	notFullCond  *cond

	head int
	tail int
	size int
}

func NewConcurrentBlockingQueue[T any](size int) *ConcurrentBlockingQueueV1[T] {
	m := &sync.Mutex{}
	return &ConcurrentBlockingQueueV1[T]{
		data: make([]any, size),
		//notEmpty: make(chan struct{}, 1),
		//notFull:  make(chan struct{}, 1),
		lock:         m,
		notEmptyCond: &cond{Cond: sync.NewCond(m)},
		notFullCond:  &cond{Cond: sync.NewCond(m)},
	}
}

func (c *ConcurrentBlockingQueueV1[T]) Enqueue(ctx context.Context, t T) error {
	c.lock.Lock()
	for c.isFull() {
		if err := c.notFullCond.WaitWithTimeout(ctx); err != nil {
			return err
		}
	}
	c.data[c.tail] = t
	c.tail++
	if c.tail == len(c.data) {
		c.tail = 0
	}
	c.size++
	// 如果有人在等待有数据 则应该唤醒这个人 但是不能阻塞我自己
	c.notEmptyCond.Signal()
	c.lock.Unlock()
	return nil
}

func (c *ConcurrentBlockingQueueV1[T]) Dequeue(ctx context.Context) (t T, err error) {
	c.lock.Lock()
	for c.isEmpty() {
		if err = c.notEmptyCond.WaitWithTimeout(ctx); err != nil {
			return
		}

	}
	t = c.data[c.head]
	c.head++
	if c.head == len(c.data) {
		c.head = 0
	}
	c.size--
	// 如果有人在等待少个数据 则应该唤醒这个人 但是不能阻塞我自己
	c.notFullCond.Signal()
	c.lock.Unlock()
	return
}

func (c *ConcurrentBlockingQueueV1[T]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.isEmpty()
}

// 封装一层 避免内部方法使用时重复加锁
func (c *ConcurrentBlockingQueueV1[T]) isEmpty() bool {
	return c.size == 0
}

func (c *ConcurrentBlockingQueueV1[T]) IsFull() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.isFull()
}

// 封装一层 避免内部方法使用时重复加锁
func (c *ConcurrentBlockingQueueV1[T]) isFull() bool {
	return c.size == len(c.data)
}

func (c *ConcurrentBlockingQueueV1[T]) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.size
}

/**
sync.Cond.WithTimeout
*/
// solution 1
type cond struct {
	*sync.Cond
}

func (c *cond) WaitWithTimeout(ctx context.Context) error {
	ch := make(chan struct{})
	go func() {
		c.Cond.Wait()
		select {
		case ch <- struct{}{}:
		default:
			c.Cond.Signal()
			c.Cond.L.Unlock()
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}

func (c *cond) Single() {
	c.Cond.Signal()
}

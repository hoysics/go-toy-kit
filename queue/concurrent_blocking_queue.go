package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

var _ Queue[any] = &ConcurrentBlockingQueue[any]{}

type ConcurrentBlockingQueue[T any] struct {
	lock *sync.Mutex
	data []T

	notEmptyCond *Cond
	notFullCond  *Cond

	head    int
	tail    int
	size    int
	maxSize int

	zero T
}

func NewConcurrentBlockingQueue[T any](maxSize int) *ConcurrentBlockingQueue[T] {
	m := &sync.Mutex{}
	return &ConcurrentBlockingQueue[T]{
		data:         make([]T, maxSize),
		maxSize:      maxSize,
		lock:         m,
		notEmptyCond: NewCond(m),
		notFullCond:  NewCond(m),
	}
}

func (c *ConcurrentBlockingQueue[T]) Get(index int) (T, error) {
	//TODO 支持随机访问
	var t T
	return t, nil
}

func (c *ConcurrentBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	c.lock.Lock()
	for c.isFull() {
		if err := c.notFullCond.WaitWithTimeout(ctx); err != nil {
			return err
		}
	}
	c.data[c.tail] = t
	c.tail++
	if c.tail == c.maxSize {
		c.tail = 0
	}
	c.size++
	// 如果有人在等待有数据 则应该唤醒这个人 但是不能阻塞我自己
	c.notEmptyCond.Broadcast()
	c.lock.Unlock()
	return nil
}

func (c *ConcurrentBlockingQueue[T]) Dequeue(ctx context.Context) (t T, err error) {
	c.lock.Lock()
	for c.isEmpty() {
		if err = c.notEmptyCond.WaitWithTimeout(ctx); err != nil {
			return
		}

	}
	t = c.data[c.head]
	c.data[c.head] = c.zero
	c.head++
	if c.head == c.maxSize {
		c.head = 0
	}
	c.size--
	// 如果有人在等待少个数据 则应该唤醒这个人 但是不能阻塞我自己
	c.notFullCond.Broadcast()
	c.lock.Unlock()
	return
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.isEmpty()
}

// 封装一层 避免内部方法使用时重复加锁
func (c *ConcurrentBlockingQueue[T]) isEmpty() bool {
	return c.size == 0
}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.isFull()
}

// 封装一层 避免内部方法使用时重复加锁
func (c *ConcurrentBlockingQueue[T]) isFull() bool {
	return c.maxSize == c.size
}

func (c *ConcurrentBlockingQueue[T]) Len() uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return uint64(c.maxSize)
}

// Cond
// solution 2 Broadcast
type Cond struct {
	L sync.Locker
	n unsafe.Pointer
}

func NewCond(l sync.Locker) *Cond {
	n := make(chan struct{})
	return &Cond{L: l, n: unsafe.Pointer(&n)}
}

func (c *Cond) Wait() {
	n := c.NotifyChan()
	c.L.Unlock()
	<-n
	c.L.Lock()
}

func (c *Cond) WaitWithTimeout(ctx context.Context) error {
	n := c.NotifyChan()
	c.L.Unlock()
	select {
	case <-ctx.Done():
		c.L.Lock()
		return ctx.Err()
	case <-n:
		c.L.Lock()
		return nil
	}
}

func (c *Cond) NotifyChan() <-chan struct{} {
	ptr := atomic.LoadPointer(&c.n)
	return *((*chan struct{})(ptr))
}

func (c *Cond) Broadcast() {
	n := make(chan struct{})
	ptrOld := atomic.SwapPointer(&c.n, unsafe.Pointer(&n))
	close(*(*chan struct{})(ptrOld))
}

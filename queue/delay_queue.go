package queue

import (
	"context"
	"sync"
	"time"
)

//type delayNode struct {
//	time time.Time
//}
//
//func compareDelay(src *delayNode, dst *delayNode) int {
//	if src.time.Equal(dst.time) {
//		return 0
//	}
//	if src.time.Before(dst.time) {
//		return -1
//	}
//	return 1
//}

type Delayable interface {
	Delay() time.Duration
}

type DelayQueue[T Delayable] struct {
	q *PriorityQueue[T]

	notEmptyCond    *Cond
	notFullCond     *Cond
	notEarliestCond *Cond

	mu sync.Mutex
}

func NewDelayQueue[T Delayable](capacity int) *DelayQueue[T] {
	return &DelayQueue[T]{
		q: NewPriorityQueue[T](10, func(src T, dst T) int {
			srcDelay := src.Delay()
			dstDelay := dst.Delay()
			if srcDelay < dstDelay {
				return -1
			}
			if srcDelay == dstDelay {
				return 0
			}
			return 1
		}),
	}
}

// Enqueue 入队和并发阻塞队列一样
func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
	d.mu.Lock()
	if ctx.Err() != nil {
		d.mu.Unlock()
		return ctx.Err()
	}
	for d.q.isFull() {
		if err := d.notFullCond.WaitWithTimeout(ctx); err != nil {
			d.mu.Unlock()
			return err
		}
	}
	// 如果如入队后的元素， 过期时间更短，那么就要唤醒出队
	// 或者一点都不管 直接唤醒出队
	top, err := d.q.Peek()
	if err != nil {
		d.mu.Unlock()
		return err
	}
	topTime := time.Now().Add(top.Delay())
	tTime := time.Now().Add(t.Delay())
	if topTime.After(tTime) {
		d.notEarliestCond.Broadcast()
	}
	err = d.q.Enqueue(t)
	d.notEmptyCond.Broadcast()
	d.mu.Unlock()
	return err
}

// Dequeue 出队
//  1. Delay()返回<=0的时候才能出队
//  2. 如果队首的Delay()(假设300ms)>0 要sleep 等待Delay()降下去
//  3. 如果正在sleep的过程 有新元素入队
//     并且Delay(假设200ms) 比正在sleep的时间还要短 需要调整sleep的时间
//  4. 如果sleep的时间还没到 就超时了 则直接返回
//
// sleep本质上是阻塞（可以用time.Sleep 也可以用channel）
func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	var sleep *time.Timer
	for {
		if ctx.Err() != nil {
			var t T
			return t, ctx.Err()
		}
		d.mu.Lock()
		if ctx.Err() != nil {
			var t T
			return t, ctx.Err()
		}
		if d.q.isEmpty() {
			if err := d.notEmptyCond.WaitWithTimeout(ctx); err != nil {
				d.mu.Unlock()
				var t T
				return t, err
			}
		}
		//等待
		head, err := d.q.Peek()
		if err != nil {
			d.mu.Unlock()
			var t T
			return t, err
		}
		if head.Delay() <= 0 {
			break
		}
		n := d.notEarliestCond.NotifyChan()
		d.mu.Unlock()
		if sleep == nil {
			sleep = time.NewTimer(head.Delay())
		} else {
			sleep.Reset(head.Delay())
		}
		select {
		case <-sleep.C:
			sleep.Stop()
			//如果有两个同时在等 那么这里其中一个会抢走，另外一个就获取不到元素了
			//d.mu.Lock()
			continue
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		case <-n:
			continue
		}
	}
	t, err := d.q.Dequeue()
	d.notFullCond.Broadcast()
	d.mu.Unlock()
	return t, err
}

func (d *DelayQueue[T]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (d *DelayQueue[T]) IsFull() bool {
	//TODO implement me
	panic("implement me")
}

func (d *DelayQueue[T]) Len() uint64 {
	//TODO implement me
	panic("implement me")
}

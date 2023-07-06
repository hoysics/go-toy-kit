package queue

import (
	"context"
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
func (d DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
}

// Dequeue 出队
//  1. Delay()返回<=0的时候才能出队
//  2. 如果队首的Delay()(假设300ms)>0 要sleep 等待Delay()降下去
//  3. 如果正在sleep的过程 有新元素入队
//     并且Delay(假设200ms) 比正在sleep的时间还要短 需要调整sleep的时间
//  4. 如果sleep的时间还没到 就超时了 则直接返回
//
// sleep本质上是阻塞（可以用time.Sleep 也可以用channel）
func (d DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (d DelayQueue[T]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (d DelayQueue[T]) IsFull() bool {
	//TODO implement me
	panic("implement me")
}

func (d DelayQueue[T]) Len() uint64 {
	//TODO implement me
	panic("implement me")
}

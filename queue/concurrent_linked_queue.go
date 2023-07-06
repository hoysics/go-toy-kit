package queue

import (
	"context"
	"sync/atomic"
	"unsafe"
)

type node[T any] struct {
	next unsafe.Pointer
	val  T
}

type ConcurrentLinkedQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func (c *ConcurrentLinkedQueue[T]) Enqueue(ctx context.Context, t T) error {
	newNode := &node[T]{
		val: t,
	}
	newPtr := unsafe.Pointer(newNode)
	//for {
	//	//方法1. 先换tail.next
	//	tailPtr := atomic.LoadPointer(&c.tail)
	//	tail := (*node[T])(tailPtr)
	//	next := atomic.LoadPointer(&tail.next)
	//	if next != nil {
	//		continue
	//	}
	//	//详情见ekit
	//	if atomic.CompareAndSwapPointer(&tail.next, next, newPtr) {
	//		atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr)
	//		return nil
	//	}
	//}
	for {
		select {
		case <-ctx.Done():
		default:
		}
		//方法2. 先换tail
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		if atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr) {
			// tail.next->c.tail
			atomic.StorePointer(&tail.next, newPtr)
			return nil
		}
	}
}

func (c *ConcurrentLinkedQueue[T]) Dequeue(ctx context.Context) (T, error) {
	for {
		if ctx.Err() != nil {
			var t T
			return t, ctx.Err()
		}
		headPtr := atomic.LoadPointer(&c.head)
		head := (*node[T])(headPtr)
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		if head == tail {
			var t T
			return t, nil // TODO 空队列
		}
		nextPtr := atomic.LoadPointer(&head.next)
		if !atomic.CompareAndSwapPointer(&c.head, headPtr, nextPtr) {
			continue
		}
		next := (*node[T])(nextPtr)
		//head.next = nil
		return next.val, nil
	}
}

func (c *ConcurrentLinkedQueue[T]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ConcurrentLinkedQueue[T]) IsFull() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ConcurrentLinkedQueue[T]) Len() uint64 {
	//TODO implement me
	panic("implement me")
}

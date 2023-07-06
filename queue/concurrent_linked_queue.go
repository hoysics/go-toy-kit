package queue

import (
	"context"
)

type node struct {
	next *node
	val  any
}

type ConcurrentLinkedQueue[T any] struct {
	head *node
	tail *node
}

func (c *ConcurrentLinkedQueue[T]) Enqueue(ctx context.Context, t T) error {
	c.tail.next = &node{
		val: t,
	}
	c.tail = c.tail.next
	return nil
}

func (c *ConcurrentLinkedQueue[T]) Dequeue(ctx context.Context) (T, error) {
	t := c.head
	c.head = c.head.next
	t.next = nil
	return t.val, nil
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

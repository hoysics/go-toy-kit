package queue

import (
	"context"
)

type Queue[T any] interface {
	Enqueue(context.Context, T) error
	Dequeue(context.Context) (T, error)
	IsEmpty() bool
	IsFull() bool
	Len() int
}

package queue

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentBlockingQueue(t *testing.T) {
	// 只能确保没有死锁
	q := NewConcurrentBlockingQueue[int](10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(30)
	//Enqueue
	for i := 0; i < 20; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				//并发状态下 没有办法校验中间结果
				ctx, cancel := context.WithTimeout(ctx, time.Second)
				err := q.Enqueue(ctx, rand.Int())
				cancel()
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}

	//Dequeue
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				ctx, cancel := context.WithTimeout(ctx, time.Second)
				_, err := q.Dequeue(ctx)
				cancel()
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}
	//校验
	wg.Wait()
}

//	func TestConcurrentBlockingQueue_Dequeue(t *testing.T) {
//		type args struct {
//			ctx context.Context
//		}
//		type testCase[T any] struct {
//			name    string
//			c       ConcurrentBlockingQueue[T]
//			args    args
//			wantT   T
//			wantErr bool
//		}
//		tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				gotT, err := tt.c.Dequeue(tt.args.ctx)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("Dequeue() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(gotT, tt.wantT) {
//					t.Errorf("Dequeue() gotT = %v, want %v", gotT, tt.wantT)
//				}
//			})
//		}
//	}
func TestConcurrentBlockingQueue_Enqueue(t *testing.T) {
	type testCase[T any] struct {
		name string
		c    *ConcurrentBlockingQueue[int]
		//ctx     context.Context
		timeout time.Duration
		input   int

		wantErr error
		data    []int
	}
	tests := []testCase[int]{
		{
			name: "enqueue",
			c:    NewConcurrentBlockingQueue[int](10),
			//ctx:   context.Background(),
			timeout: time.Minute,
			input:   1,
			data:    []int{1},
			wantErr: nil,
		},
		{
			name: "blocking",
			c: func() *ConcurrentBlockingQueue[int] {
				res := NewConcurrentBlockingQueue[int](2)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := res.Enqueue(ctx, 1)
				require.NoError(t, err)
				err = res.Enqueue(ctx, 2)
				require.NoError(t, err)
				return res
			}(),
			timeout: time.Second,
			data:    []int{1, 2},
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			if err := tt.c.Enqueue(ctx, tt.input); !errors.Is(err, tt.wantErr) {
				t.Errorf("Enqueue() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.c.data, tt.data)
		})
	}
}

//
//func TestConcurrentBlockingQueue_IsEmpty(t *testing.T) {
//	type testCase[T any] struct {
//		name string
//		c    ConcurrentBlockingQueue[T]
//		want bool
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.IsEmpty(); got != tt.want {
//				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestConcurrentBlockingQueue_IsFull(t *testing.T) {
//	type testCase[T any] struct {
//		name string
//		c    ConcurrentBlockingQueue[T]
//		want bool
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.IsFull(); got != tt.want {
//				t.Errorf("IsFull() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestConcurrentBlockingQueue_Len(t *testing.T) {
//	type testCase[T any] struct {
//		name string
//		c    ConcurrentBlockingQueue[T]
//		want int
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.Len(); got != tt.want {
//				t.Errorf("Len() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestConcurrentBlockingQueue_isEmpty(t *testing.T) {
//	type testCase[T any] struct {
//		name string
//		c    ConcurrentBlockingQueue[T]
//		want bool
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.isEmpty(); got != tt.want {
//				t.Errorf("isEmpty() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestConcurrentBlockingQueue_isFull(t *testing.T) {
//	type testCase[T any] struct {
//		name string
//		c    ConcurrentBlockingQueue[T]
//		want bool
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.isFull(); got != tt.want {
//				t.Errorf("isFull() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestCond_Broadcast(t *testing.T) {
//	type fields struct {
//		L sync.Locker
//		n unsafe.Pointer
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Cond{
//				L: tt.fields.L,
//				n: tt.fields.n,
//			}
//			c.Broadcast()
//		})
//	}
//}
//
//func TestCond_NotifyChan(t *testing.T) {
//	type fields struct {
//		L sync.Locker
//		n unsafe.Pointer
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   <-chan struct{}
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Cond{
//				L: tt.fields.L,
//				n: tt.fields.n,
//			}
//			if got := c.NotifyChan(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NotifyChan() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestCond_Wait(t *testing.T) {
//	type fields struct {
//		L sync.Locker
//		n unsafe.Pointer
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Cond{
//				L: tt.fields.L,
//				n: tt.fields.n,
//			}
//			c.Wait()
//		})
//	}
//}
//
//func TestCond_WaitWithTimeout(t *testing.T) {
//	type fields struct {
//		L sync.Locker
//		n unsafe.Pointer
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Cond{
//				L: tt.fields.L,
//				n: tt.fields.n,
//			}
//			if err := c.WaitWithTimeout(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("WaitWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestNewConcurrentBlockingQueue(t *testing.T) {
//	type args struct {
//		size int
//	}
//	type testCase[T any] struct {
//		name string
//		args args
//		want *ConcurrentBlockingQueue[T]
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewConcurrentBlockingQueue(tt.args.size); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewConcurrentBlockingQueue() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestNewCond(t *testing.T) {
//	type args struct {
//		l sync.Locker
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Cond
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewCond(tt.args.l); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewCond() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

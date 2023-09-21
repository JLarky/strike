package promise

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
)

var reqCounter uint64

type Promise[T any] struct {
	PromiseId string
	readyCh   chan bool
	value     T
	ctx       context.Context
}

type ctxKey int

const chunkQueueKey ctxKey = 0

type PromiseResult struct {
	PromiseId string
	Result    any
}

type chunkQueue struct {
	chunkChannel chan any
	workerGroup  *sync.WaitGroup
}

func WithContext(ctx context.Context) (context.Context, func() chan any) {
	chunkCh := make(chan any)
	wg := sync.WaitGroup{}

	chunkQueue := chunkQueue{
		chunkChannel: chunkCh,
		workerGroup:  &wg,
	}

	getChunkCh := func() chan any {
		go func() {
			workGroupDone := make(chan bool)
			go func() {
				wg.Wait()
				fmt.Println("WorkerGroup done")
				workGroupDone <- true
			}()
			select {
			case <-ctx.Done():
			case <-workGroupDone:
			}
			close(chunkCh)
		}()

		return chunkCh
	}

	return context.WithValue(ctx, chunkQueueKey, &chunkQueue), getChunkCh
}

func FromContext(ctx context.Context) (*chunkQueue, bool) {
	queue, ok := ctx.Value(chunkQueueKey).(*chunkQueue)
	return queue, ok
}

func NewPromise[T any](ctx context.Context) Promise[T] {
	id := atomic.AddUint64(&reqCounter, 1)
	c := make(chan bool)
	return Promise[T]{
		PromiseId: fmt.Sprintf("%d", id),
		readyCh:   c,
		ctx:       ctx,
	}
}

func (p *Promise[T]) Resolve(value T) {
	p.value = value
	close(p.readyCh)
}

func (c Promise[T]) String() string {
	type Alias Promise[T]
	return fmt.Sprintf("Promise(%v)", Alias(c))
}

func (p Promise[T]) MarshalJSON() ([]byte, error) {
	p.MarkSent()
	fmt.Println("Promise MarshalJSON", p)
	return json.Marshal(map[string]any{
		"$strike": "promise",
		"id":      p.PromiseId,
	})
}

func (p PromiseResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"$strike": "promise-result",
		"id":      p.PromiseId,
		"result":  p.Result,
	})
}

func (p *Promise[T]) ResolveAsync(valueGen func() T) {
	queue, ok := FromContext(p.ctx)
	if !ok {
		panic("failed to get context")
	}

	queue.workerGroup.Add(1)

	go func() {
		defer queue.workerGroup.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%T\n", p)
				fmt.Printf("%T\n", valueGen)
				fmt.Println("Recovered in ResolveAsync", r)
			}
		}()
		// exit early if the context is done
		select {
		case <-p.ctx.Done():
			return
		default:
			// continue
		}

		ch := make(chan PromiseResult)
		go func() {
			v := valueGen()
			p.Resolve(v)
			ch <- PromiseResult{
				PromiseId: p.PromiseId,
				Result:    v,
			}
		}()
		// chunkChannel is going to be close if the context is done
		select {
		case <-p.ctx.Done():
		case chunk := <-ch:
			queue.chunkChannel <- chunk
		}
	}()
}

func (p *Promise[T]) Then() T {
	<-p.readyCh
	return p.value
}

func (p *Promise[T]) MarkSent() {
	fmt.Println("MarkSent", p.ctx.Value(chunkQueueKey))
	p.readyCh = nil
}

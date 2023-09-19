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
	ready     chan bool
	value     T
	ctx       context.Context
}

type test string

type Task struct {
	ID     string
	Result any
}

type MyStuff struct {
	TaskChannel chan Task
	WorkerGroup *sync.WaitGroup
}

func WithContext(ctx context.Context) (context.Context, func() chan Task) {
	taskChannel := make(chan Task)
	wg := sync.WaitGroup{}

	myStuff := MyStuff{
		TaskChannel: taskChannel,
		WorkerGroup: &wg,
	}

	getTaskCh := func() chan Task {
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
			close(taskChannel)
		}()

		return taskChannel
	}

	return context.WithValue(ctx, test("promise"), &myStuff), getTaskCh
}

func FromContext(ctx context.Context) (*MyStuff, bool) {
	stuff, ok := ctx.Value(test("promise")).(*MyStuff)
	return stuff, ok
}

func NewPromise[T any](ctx context.Context) Promise[T] {
	id := atomic.AddUint64(&reqCounter, 1)
	c := make(chan bool)
	return Promise[T]{
		PromiseId: fmt.Sprintf("%d", id),
		ready:     c,
		ctx:       ctx,
	}
}

func (p *Promise[T]) Resolve(value T) {
	p.value = value
	close(p.ready)
}

func (p Promise[T]) MarshalJSON() ([]byte, error) {
	p.MarkSent()
	fmt.Println("Promise MarshalJSON", p)
	return json.Marshal(map[string]any{
		"$strike": "promise",
		"id":      p.PromiseId,
	})
}

func (p *Promise[T]) ResolveAsync(valueGen func() T) {
	myStuff, ok := FromContext(p.ctx)
	if !ok {
		panic("failed to get context")
	}

	myStuff.WorkerGroup.Add(1)

	go func() {
		defer myStuff.WorkerGroup.Done()
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

		ch := make(chan Task)
		go func() {
			v := valueGen()
			p.Resolve(v)
			ch <- Task{
				ID:     p.PromiseId,
				Result: v,
			}
		}()
		// TaskChannel is going to be close if the context is done
		select {
		case <-p.ctx.Done():
		case task := <-ch:
			myStuff.TaskChannel <- task
		}
	}()
}

func (p *Promise[T]) Then() T {
	<-p.ready
	return p.value
}

func (p *Promise[T]) MarkSent() {
	fmt.Println("MarkSent", p.ctx.Value(test("promise")))
	p.ready = nil
}

func (c Promise[T]) String() string {
	type Alias Promise[T]
	return fmt.Sprintf("Promise(%v)", Alias(c))
}

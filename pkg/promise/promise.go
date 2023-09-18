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
	ID     int
	Result string
}

type MyStuff struct {
	TaskChannel chan Task
	WorkerGroup *sync.WaitGroup
}

func WithContext(ctx context.Context) (context.Context, chan Task, *sync.WaitGroup) {

	taskChannel := make(chan Task)
	wg := sync.WaitGroup{}

	myStuff := MyStuff{
		TaskChannel: taskChannel,
		WorkerGroup: &wg,
	}

	return context.WithValue(ctx, test("promise"), &myStuff), taskChannel, &wg
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

	fmt.Println("ResolveAsync", myStuff)

	go func() {
		defer myStuff.WorkerGroup.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%T\n", p)
				fmt.Printf("%T\n", valueGen)
				fmt.Println("Recovered in ResolveAsync", r)
			}
		}()
		fmt.Println("ResolveAsync before")
		v := valueGen()
		myStuff.TaskChannel <- Task{
			ID:     11,
			Result: fmt.Sprintf("New component %v", v),
		}
		fmt.Println("ResolveAsync", v)
		p.Resolve(v)
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

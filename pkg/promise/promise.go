package promise

import (
	"fmt"
	"sync/atomic"
)

var reqCounter uint64

type Promise[T any] struct {
	PromiseId string
	ready     chan bool
	value     T
}

func NewPromise[T any]() Promise[T] {
	id := atomic.AddUint64(&reqCounter, 1)
	c := make(chan bool)
	return Promise[T]{
		PromiseId: fmt.Sprintf("%d", id),
		ready:     c,
	}
}

func (p *Promise[T]) Resolve(value T) {
	p.value = value
	close(p.ready)
}

func (p *Promise[T]) ResolveAsync(valueGen func() T) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%T\n", p)
				fmt.Printf("%T\n", valueGen)
				fmt.Println("Recovered in ResolveAsync", r)
			}
		}()
		fmt.Println("ResolveAsync before")
		v := valueGen()
		fmt.Println("ResolveAsync", v)
		p.Resolve(v)
	}()
}

func (p *Promise[T]) Then() T {
	<-p.ready
	return p.value
}

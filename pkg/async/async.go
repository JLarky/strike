package async

import (
	"github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/promise"
)

func Async(compGen func() h.Component) promise.Promise[h.Component] {
	p := promise.NewPromise[h.Component]()
	return p
}

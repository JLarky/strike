package async

import (
	"context"

	"github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/promise"
)

func Async(ctx context.Context, compGen func() h.Component) promise.Promise[h.Component] {
	p := promise.NewPromise[h.Component](ctx)
	p.ResolveAsync(compGen)
	return p
}

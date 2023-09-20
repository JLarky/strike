package suspense

import (
	"context"
	"log"

	"github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/promise"
)

const suspenseComponentName = "strike-suspense"

func Suspense(comp h.Component) h.Component {
	comp.Tag_type = suspenseComponentName

	// validate ctx prop; remove if invalid
	if ctx, ok := comp.Props["ctx"].(context.Context); !ok {
		log.Printf("property ctx should be passed to suspense component (%v)", comp)
	} else {
		if _, ok := promise.FromContext(ctx); !ok {
			log.Printf("property ctx should be created with promise.WithContext (%v)", comp)
			delete(comp.Props, "ctx")
		}
	}

	ctx, canStream := comp.Props["ctx"].(context.Context)
	if !canStream {
		comp.Props["cantStream"] = true
		comp.Props["fallback"] = nil
	}
	children := comp.Props["children"].([]any)
	for k, v := range children {
		switch v := v.(type) {
		case func() h.Component:
			if canStream {
				p := promise.NewPromise[h.Component](ctx)
				p.ResolveAsync(v)
				children[k] = p
			} else {
				children[k] = v()
			}
		}
	}

	return comp
}

func IsSuspense(comp h.Component) bool {
	return comp.Tag_type == suspenseComponentName
}

func CanStream(comp h.Component) bool {
	return nil != comp.Props["ctx"]
}

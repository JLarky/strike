package suspense

import (
	"github.com/JLarky/strike/pkg/h"
)

const suspenseComponentName = "strike-suspense"

func Suspense(comp h.Component) h.Component {
	comp.Tag_type = suspenseComponentName
	return comp
}

func IsSuspense(comp h.Component) bool {
	return comp.Tag_type == suspenseComponentName
}

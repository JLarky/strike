package island

import (
	"github.com/JLarky/strike/pkg/h"
)

const islandComponentName = "strike-island"

func Island(comp h.Component) h.Component {
	comp.Tag_type = islandComponentName
	return comp
}

func IsIsland(comp h.Component) bool {
	return comp.Tag_type == islandComponentName
}

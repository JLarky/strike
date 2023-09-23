package action

import (
	"log"

	"github.com/JLarky/strike/pkg/h"
)

// for fallback if form was not hydrated

func HiddenInput(comp h.Component) h.Component {
	comp.Tag_type = "input"
	comp.Props["type"] = "hidden"

	if comp.Props["name"] == nil {
		log.Println("name is required for hidden input")
	}
	if val, ok := comp.Props["name"].(Action); ok {
		comp.Props["name"] = "$ACTION_ID_" + val.Id
	} else {
		log.Println("prop `name` must be an action")
	}

	return comp
}
